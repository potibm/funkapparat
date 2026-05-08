package formatters

import (
	"fmt"
	"time"

	"github.com/gorilla/feeds"
	"github.com/potibm/funkapparat/internal/app/domain"
)

type FeedFormat string

const (
	RSSFormat  FeedFormat = "rss"
	JSONFormat FeedFormat = "json"
	AtomFormat FeedFormat = "atom"
)

type FeedFormatter struct {
	Feedtype        FeedFormat
	FeedTitle       string
	FeedLink        string
	FeedDescription string
	AuthorName      string
	AuthorEmail     string
}

func NewFeedFormatter(feedtype FeedFormat, title, description, link, authorName, authorEmail string) *FeedFormatter {
	if title == "" {
		title = "Party Announcements"
	}

	if link == "" {
		link = "https://news.scene.org"
	}

	return &FeedFormatter{
		Feedtype:        feedtype,
		FeedTitle:       title,
		FeedLink:        link,
		FeedDescription: description,
		AuthorName:      authorName,
		AuthorEmail:     authorEmail,
	}
}

func (f *FeedFormatter) Extension() string {
	switch f.Feedtype {
	case JSONFormat:
		return ".json"
	case AtomFormat:
		return ".atom"
	default:
		return ".xml"
	}
}

func (f *FeedFormatter) Format(announcements domain.AnnouncementList) ([]byte, error) {
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

	resultString := ""

	var err error

	switch f.Feedtype {
	case JSONFormat:
		resultString, err = feed.ToJSON()
	case AtomFormat:
		resultString, err = feed.ToAtom()
	default:
		resultString, err = feed.ToRss()
	}

	if err != nil {
		return nil, err
	}

	return []byte(resultString), nil
}
