package repository

import (
	"context"

	"github.com/potibm/funkapparat/internal/app/domain"
)

type AnnouncementListFilters struct {
	Query  *string
	ID     *int64
	Hidden *bool
}

type AnnouncementListParams struct {
	Offset int
	Limit  int
	Sort   string
	Order  string
}

type AnnouncementRepository interface {
	Save(ctx context.Context, announcement *domain.Announcement) error
	Delete(ctx context.Context, id int64) error
	List(
		ctx context.Context,
		params AnnouncementListParams,
		filters AnnouncementListFilters,
	) ([]domain.Announcement, int64, error)
	GetByID(ctx context.Context, id int64) (*domain.Announcement, error)
	GetAll(ctx context.Context) (domain.AnnouncementList, error)
}
