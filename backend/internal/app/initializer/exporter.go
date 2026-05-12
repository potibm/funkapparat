package initializer

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/potibm/funkapparat/internal/app/config"
	"github.com/potibm/funkapparat/internal/app/exporter"
	"github.com/potibm/funkapparat/internal/app/exporter/formatters"
	"github.com/potibm/funkapparat/internal/app/exporter/writers"
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
		var f exporter.Formatter

		switch cfg.Type {
		case "rss", "json", "atom":
			if feedConfig == nil {
				return nil, fmt.Errorf("exporter %s requires a feed config", cfg.Name)
			}

			f = formatters.NewFeedFormatter(
				formatters.FeedFormat(cfg.Type),
				feedConfig.FeedTitle,
				feedConfig.FeedDescription,
				feedConfig.FeedLink,
				feedConfig.AuthorName,
				feedConfig.AuthorEmail,
			)

			slog.Info("Configured RSS exporter", "name", cfg.Name)
		default:
			baseLog.Error("Unknown exporter type", "type", cfg.Type)

			continue
		}

		var w exporter.Writer

		switch cfg.Destination {
		case "s3":
			if s3Client == nil {
				return nil, fmt.Errorf("exporter %s requires s3, but s3client is not configured", cfg.Name)
			}

			bucket := cfg.Options["bucket"]
			if bucket == "" {
				return nil, fmt.Errorf("exporter %s: destination 's3' requires 'bucket' option", cfg.Name)
			}

			w = writers.NewS3Writer(s3Client, bucket)
		case "file":
			dir := cfg.Options["dir"]
			if dir == "" {
				return nil, fmt.Errorf("exporter %s: destination 'file' requires 'dir' option", cfg.Name)
			}

			w = &writers.FileWriter{BaseDir: dir}
		default:
			baseLog.Error("Unknown destination", "dest", cfg.Destination)

			continue
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
