package seeder

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/potibm/funkapparat/internal/app/domain"
	"github.com/potibm/funkapparat/internal/app/repository"
)

const (
	numberOfAnnouncements = 15
)

type Seeder struct {
	repo repository.AnnouncementRepository
}

func NewSeeder(repo repository.AnnouncementRepository) *Seeder {
	return &Seeder{
		repo: repo,
	}
}

func (s *Seeder) Run() error {
	ctx := context.Background()

	slog.Info("Starting DB Seed...")

	_ = gofakeit.Seed(0)

	c := 0

	for range numberOfAnnouncements {
		var url string
		if gofakeit.Bool() { // 50/50 Chance
			url = gofakeit.URL()
		} else {
			url = ""
		}

		announcement := domain.Announcement{
			Title:       gofakeit.Sentence(5),
			Body:        generateMarkdownCodeBody(),
			IsUrgent:    false,
			ExternalURL: url,
			IsHidden:    gofakeit.Bool(),
		}

		if err := s.repo.Save(ctx, &announcement); err != nil {
			slog.Error("Failed to seed announcement", "error", err)
		} else {
			slog.Info("Seeded announcement", "id", announcement.ID)

			c++
		}
	}

	slog.Info("Seeding finished successfully", "total_announcements", c)

	return nil
}

func generateMarkdownCodeBody() string {
	var b strings.Builder

	if gofakeit.Bool() {
		b.WriteString("## " + gofakeit.Sentence(3) + "\n\n")
	}

	b.WriteString(generateMarkdownParagraph() + "\n\n")

	if gofakeit.Bool() {
		b.WriteString("### " + gofakeit.Sentence(3) + "\n\n")
		b.WriteString(generateMarkdownParagraph() + "\n\n")
	}

	if gofakeit.Bool() {
		b.WriteString("### " + gofakeit.Sentence(3) + "\n\n")

		count := gofakeit.Number(3, 5)
		for i := 0; i < count; i++ {
			fmt.Fprintf(&b, "- %s\n", gofakeit.Sentence(4))
		}

		b.WriteString("\n")
	}

	return b.String()
}

func generateMarkdownParagraph() string {
	sentenceCount := gofakeit.Number(2, 5)

	var sentences []string

	for i := 0; i < sentenceCount; i++ {
		sentence := gofakeit.Sentence(gofakeit.Number(5, 10))

		switch gofakeit.Number(0, 15) {
		case 0:
			sentence = "**" + sentence + "**"
		case 1:
			sentence = "_" + sentence + "_"
		case 2:
			sentence = fmt.Sprintf("[%s](%s)", gofakeit.Word(), gofakeit.URL())
		}

		sentences = append(sentences, sentence)
	}

	return strings.Join(sentences, " ")
}
