package domain

import (
	"bytes"
	"encoding/json"
	"html"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
)

var (
	mdConverter = goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
	)

	strictSanitizer = bluemonday.StrictPolicy()
)

type Announcement struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Body        string    `json:"body"`
	IsUrgent    bool      `json:"is_urgent"`
	ExternalURL string    `json:"external_url,omitempty"`
	IsHidden    bool      `json:"is_hidden"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AnnouncementList []*Announcement

func (a *Announcement) HTML() string {
	if a.Body == "" {
		return ""
	}

	var buf bytes.Buffer
	if err := mdConverter.Convert([]byte(a.Body), &buf); err != nil {
		return a.Body
	}

	return buf.String()
}

func (a *Announcement) PlainText() string {
	htmlContent := a.HTML()

	plainText := strictSanitizer.Sanitize(htmlContent)

	return html.UnescapeString(plainText)
}

func (a *Announcement) MarshalJSON() ([]byte, error) {
	type Alias Announcement

	return json.Marshal(&struct {
		*Alias

		BodyHTML      string `json:"body_html"`
		BodyPlainText string `json:"body_plain"`
	}{
		Alias:         (*Alias)(a),
		BodyHTML:      a.HTML(),
		BodyPlainText: a.PlainText(),
	})
}
