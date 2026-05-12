package services

import (
	"github.com/potibm/funkapparat/internal/app/domain"
	"github.com/potibm/protokolapparat/pkg/news"
)

func mapToEventPayload(entry *domain.Announcement) news.Entry {
	result := news.Entry{
		ID:          entry.ID,
		Title:       entry.Title,
		Body:        entry.Body,
		IsUrgent:    entry.IsUrgent,
		ExternalURL: entry.ExternalURL,
		IsHidden:    entry.IsHidden,
	}

	return result
}

func mapToEventListPayload(entries domain.AnnouncementList) []news.Entry {
	result := make([]news.Entry, 0, len(entries))
	for _, entry := range entries {
		result = append(result, mapToEventPayload(entry))
	}

	return result
}
