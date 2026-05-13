package initializer

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/potibm/funkapparat/internal/app/config"
	"github.com/potibm/funkapparat/internal/app/exporter"
	"github.com/potibm/funkapparat/internal/app/exporter/formatters"
	"github.com/potibm/funkapparat/internal/app/exporter/writers"
)

var (
	errUnknownType = errors.New("unknown exporter type")
	errUnknownDest = errors.New("unknown destination")
)

func BootstrapExporters(
	ctx context.Context,
	version string,
	feedConfig *config.FeedConfig,
	configs []config.ExporterConfig,
	s3Client *s3.Client,
	baseLog *slog.Logger,
) ([]exporter.Exporter, error) {
	var result []exporter.Exporter

	exporterLog := slog.With("component", "Exporter")

	for _, cfg := range configs {
		if !cfg.Enabled {
			exporterLog.Debug("Skipping disabled exporter", "name", cfg.Name)

			continue
		}

		f, err := buildFormatter(cfg, feedConfig)
		if err != nil {
			baseLog.Error("Unknown exporter type", "type", cfg.Type)

			continue
		}

		w, err := buildWriter(cfg, s3Client)
		if err != nil {
			if errors.Is(err, errUnknownDest) {
				baseLog.Error("Unknown destination", "dest", cfg.Destination)

				continue
			}

			return nil, err
		}

		ex := exporter.NewUniversalExporter(
			cfg.Name,
			cfg.Filename,
			f,
			w,
			exporterLog.With("exporter", cfg.Name),
		)
		result = append(result, ex)
	}

	return result, nil
}

func buildFormatter(
	cfg config.ExporterConfig,
	feedConfig *config.FeedConfig,
) (exporter.Formatter, error) {
	switch cfg.Type {
	case "rss", "json", "atom":
		if feedConfig == nil {
			return nil, fmt.Errorf("exporter %s requires a feed config", cfg.Name)
		}

		return formatters.NewFeedFormatter(
			formatters.FeedFormat(cfg.Type),
			feedConfig.FeedTitle,
			feedConfig.FeedDescription,
			feedConfig.FeedLink,
			feedConfig.AuthorName,
			feedConfig.AuthorEmail,
		), nil

	default:
		return nil, errUnknownType
	}
}

func buildWriter(cfg config.ExporterConfig, s3Client *s3.Client) (exporter.Writer, error) {
	switch cfg.Destination {
	case "s3":
		if s3Client == nil {
			return nil, fmt.Errorf("exporter %s requires s3, but s3client is not configured", cfg.Name)
		}

		bucket := cfg.Options["bucket"]
		if bucket == "" {
			return nil, fmt.Errorf("exporter %s: destination 's3' requires 'bucket' option", cfg.Name)
		}

		return writers.NewS3Writer(s3Client, bucket), nil
	case "file":
		dir := cfg.Options["dir"]
		if dir == "" {
			return nil, fmt.Errorf("exporter %s: destination 'file' requires 'dir' option", cfg.Name)
		}

		return &writers.FileWriter{BaseDir: dir}, nil
	default:
		return nil, errUnknownDest
	}
}
