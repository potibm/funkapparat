package gorm

import (
	"github.com/potibm/funkapparat/internal/app/domain"
)

type dbAnnouncement struct {
	GormModel

	Title       string
	Body        string
	IsUrgent    bool
	ExternalURL string
	IsHidden    bool
}

func (dbAnnouncement) TableName() string {
	return "announcements"
}

func fromDomainAnnouncement(a *domain.Announcement) *dbAnnouncement {
	return &dbAnnouncement{
		GormModel: GormModel{ID: a.ID},

		Title:       a.Title,
		Body:        a.Body,
		IsUrgent:    a.IsUrgent,
		ExternalURL: a.ExternalURL,
		IsHidden:    a.IsHidden,
	}
}

func toDomainAnnouncement(db *dbAnnouncement) *domain.Announcement {
	return &domain.Announcement{
		ID:          db.ID,
		Title:       db.Title,
		Body:        db.Body,
		IsUrgent:    db.IsUrgent,
		ExternalURL: db.ExternalURL,
		IsHidden:    db.IsHidden,
		CreatedAt:   db.CreatedAt,
		UpdatedAt:   db.UpdatedAt,
	}
}

func toDomainAnnouncementList(db *[]dbAnnouncement) domain.AnnouncementList {
	announcements := make(domain.AnnouncementList, len(*db))
	for i, dbAnnouncement := range *db {
		announcements[i] = toDomainAnnouncement(&dbAnnouncement)
	}

	return announcements
}
