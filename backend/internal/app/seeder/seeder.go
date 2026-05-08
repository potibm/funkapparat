package seeder

import (
	"log/slog"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/potibm/funkapparat/internal/app/domain"
	"github.com/potibm/funkapparat/internal/app/repository"
)

const (
	numberOfSponsorSlides     = 10
	numberOfSceneFriendSlides = 25
	numberOfNewsSlides        = 10
)

type Seeder struct {
	annoucenments []domain.Announcement
	currentID     int64
	repo          repository.AnnouncementRepository
}

func NewSeeder(repo repository.AnnouncementRepository) *Seeder {
	return &Seeder{
		annoucenments: []domain.Announcement{},
		currentID:     0,
		repo:          repo,
	}
}

func (s *Seeder) Run() error {
	slog.Info("Starting DB Purge & Seed...")

	_ = gofakeit.Seed(0)

	slog.Info("Seeding finished successfully", "total_slides", s.currentID)

	return nil
}
