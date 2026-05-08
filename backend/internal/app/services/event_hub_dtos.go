package services

import (
	"time"

	"github.com/potibm/funkapparat/internal/app/domain"
)

type ActionType string

const (
	ActionCreate ActionType = "create"
	ActionUpdate ActionType = "update"
	ActionDelete ActionType = "delete"
	ActionSync   ActionType = "sync"
)

type AnnouncementSyncEventDTO struct {
	Action    ActionType        `json:"action"`
	Timestamp int64             `json:"timestamp"`
	Payload   []AnnouncementDTO `json:"payload"`
}

type AnnouncementEventDTO struct {
	Action    ActionType      `json:"action"`
	Timestamp int64           `json:"timestamp"`
	Payload   AnnouncementDTO `json:"payload"`
}

type AnnouncementDTO struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Body        string `json:"body"`
	IsUrgent    bool   `json:"is_urgent"`
	ExternalURL string `json:"external_url,omitempty"`
	IsHidden    bool   `json:"is_hidden"`
}

func mapToEntryDTO(entry *domain.Announcement) AnnouncementDTO {
	dto := AnnouncementDTO{
		ID:          entry.ID,
		Title:       entry.Title,
		Body:        entry.Body,
		IsUrgent:    entry.IsUrgent,
		ExternalURL: entry.ExternalURL,
		IsHidden:    entry.IsHidden,
	}

	return dto
}

func mapToEventDTO(entry *domain.Announcement, action ActionType) AnnouncementEventDTO {
	return AnnouncementEventDTO{
		Action:    action,
		Timestamp: time.Now().Unix(),
		Payload:   mapToEntryDTO(entry),
	}
}

func mapToAnnouncementListDTO(entries domain.AnnouncementList) []AnnouncementDTO {
	dtos := make([]AnnouncementDTO, 0, len(entries))
	for _, entry := range entries {
		dtos = append(dtos, mapToEntryDTO(entry))
	}

	return dtos
}
