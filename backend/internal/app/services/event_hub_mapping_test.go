package services

import (
	"testing"

	"github.com/potibm/funkapparat/internal/app/domain"
	"github.com/potibm/protokolapparat/pkg/news"
	"github.com/stretchr/testify/assert"
)

func TestMapToEventPayload(t *testing.T) {
	tests := []struct {
		name     string
		entry    *domain.Announcement
		expected news.Entry
	}{
		{
			name: "maps all fields correctly",
			entry: &domain.Announcement{
				ID:          42,
				Title:       "Test Title",
				Body:        "Test Body",
				IsUrgent:    true,
				ExternalURL: "https://example.com",
				IsHidden:    true,
			},
			expected: news.Entry{
				ID:          42,
				Title:       "Test Title",
				Body:        "Test Body",
				IsUrgent:    true,
				ExternalURL: "https://example.com",
				IsHidden:    true,
			},
		},
		{
			name: "maps zero values correctly",
			entry: &domain.Announcement{
				ID:       0,
				Title:    "",
				Body:     "",
				IsUrgent: false,
				IsHidden: false,
			},
			expected: news.Entry{
				ID:       0,
				Title:    "",
				Body:     "",
				IsUrgent: false,
				IsHidden: false,
			},
		},
		{
			name: "ignores fields not present in news.Entry",
			entry: &domain.Announcement{
				ID:    99,
				Title: "Another Title",
				Body:  "Another Body",
			},
			expected: news.Entry{
				ID:    99,
				Title: "Another Title",
				Body:  "Another Body",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapToEventPayload(tt.entry)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMapToEventListPayload(t *testing.T) {
	tests := []struct {
		name     string
		entries  domain.AnnouncementList
		expected []news.Entry
	}{
		{
			name: "maps multiple entries",
			entries: domain.AnnouncementList{
				{ID: 1, Title: "First", Body: "Body 1"},
				{ID: 2, Title: "Second", Body: "Body 2", IsUrgent: true},
				{ID: 3, Title: "Third", Body: "Body 3", ExternalURL: "https://third.com"},
			},
			expected: []news.Entry{
				{ID: 1, Title: "First", Body: "Body 1"},
				{ID: 2, Title: "Second", Body: "Body 2", IsUrgent: true},
				{ID: 3, Title: "Third", Body: "Body 3", ExternalURL: "https://third.com"},
			},
		},
		{
			name:     "returns empty slice for nil input",
			entries:  nil,
			expected: []news.Entry{},
		},
		{
			name:     "returns empty slice for empty input",
			entries:  domain.AnnouncementList{},
			expected: []news.Entry{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapToEventListPayload(tt.entries)
			assert.Equal(t, tt.expected, result)
		})
	}
}
