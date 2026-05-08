package exporter

import (
	"context"
	"log/slog"

	"github.com/potibm/funkapparat/internal/app/domain"
)

type LogExporter struct{}

func NewLogExporter() *LogExporter {
	return &LogExporter{}
}

func (e *LogExporter) Name() string {
	return "LogExporter"
}

func (e *LogExporter) Export(ctx context.Context, entries domain.AnnouncementList) error {
	slog.Info("Exporting announcements to log", "count", len(entries))

	for _, entry := range entries {
		slog.Debug("Annoucement Entry",
			"id", entry.ID,
			"title", entry.Title,
		)
	}

	return nil
}
