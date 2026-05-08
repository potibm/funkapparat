package formatters

import (
	"fmt"
	"time"

	"github.com/gorilla/feeds"
	"github.com/potibm/funkapparat/internal/app/domain"
)

type AnnouncementRssFormatter struct {
	FeedTitle       string
	FeedLink        string
	FeedDescription string
	AuthorName      string
	AuthorEmail     string
}

func NewAnnouncementRssFormatter(title, description, link, authorName, authorEmail string) *AnnouncementRssFormatter {
	if title == "" {
		title = "Party Announcements"
	}

	if link == "" {
		link = "https://news.scene.org"
	}

	return &AnnouncementRssFormatter{
		FeedTitle:       title,
		FeedLink:        link,
		FeedDescription: description,
		AuthorName:      authorName,
		AuthorEmail:     authorEmail,
	}
}

func (f *AnnouncementRssFormatter) Extension() string {
	return ".xml"
}

func (f *AnnouncementRssFormatter) Format(announcements domain.AnnouncementList) ([]byte, error) {
	feed := &feeds.Feed{
		Title:       f.FeedTitle,
		Link:        &feeds.Link{Href: f.FeedLink},
		Description: f.FeedDescription,
		Author:      &feeds.Author{Name: f.AuthorName, Email: f.AuthorEmail},
		Created:     time.Now(),
	}

	for _, ann := range announcements {
		if ann.IsHidden {
			continue
		}

		uid := fmt.Sprintf("announcement-%d@%s", ann.ID, f.FeedLink)

		itemLink := f.FeedLink
		if ann.ExternalURL != "" {
			itemLink = ann.ExternalURL
		}

		item := &feeds.Item{
			Id:          uid,
			IsPermaLink: "false",
			Title:       ann.Title,
			Link:        &feeds.Link{Href: itemLink},
			Description: ann.HTML(),
			Created:     ann.CreatedAt,
			Updated:     ann.UpdatedAt,
		}

		feed.Items = append(feed.Items, item)
	}

	rssString, err := feed.ToRss()
	if err != nil {
		return nil, err
	}

	return []byte(rssString), nil
}
