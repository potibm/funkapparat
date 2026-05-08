package gorm

import (
	"context"
	"fmt"

	"github.com/potibm/funkapparat/internal/app/domain"
	"github.com/potibm/funkapparat/internal/app/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type announcementRepository struct {
	db *gorm.DB
}

func (s *Store) NewAnnoucementRepository() repository.AnnouncementRepository {
	return NewAnnoucementRepository(s.db)
}

func NewAnnoucementRepository(db *gorm.DB) repository.AnnouncementRepository {
	return &announcementRepository{db: db}
}

func (r *announcementRepository) Save(ctx context.Context, announcement *domain.Announcement) error {
	dbObj := fromDomainAnnouncement(announcement)

	err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(dbObj).Error
	if err == nil {
		announcement.ID = dbObj.ID
	}

	return err
}

func (r *announcementRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&dbAnnouncement{}, id).Error
}

func (r *announcementRepository) List(
	ctx context.Context,
	params repository.AnnouncementListParams,
	filters repository.AnnouncementListFilters,
) ([]domain.Announcement, int64, error) {
	var (
		dbEntries []dbAnnouncement
		total     int64
	)

	query := r.db.WithContext(ctx).Model(&dbAnnouncement{})

	query = r.applyFilters(query, filters)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	query = query.Order(clause.OrderByColumn{
		Column: clause.Column{Name: params.Sort},
		Desc:   params.Order == "DESC",
	})

	err = query.Offset(params.Offset).
		Limit(params.Limit).
		Find(&dbEntries).
		Error
	if err != nil {
		return nil, 0, err
	}

	entries := make([]domain.Announcement, len(dbEntries))
	for i, dbEntry := range dbEntries {
		entries[i] = *toDomainAnnouncement(&dbEntry)
	}

	return entries, total, nil
}

func (r *announcementRepository) GetByID(ctx context.Context, id int64) (*domain.Announcement, error) {
	var dbEntry dbAnnouncement

	err := r.db.WithContext(ctx).First(&dbEntry, id).Error
	if err != nil {
		return nil, err
	}

	entry := toDomainAnnouncement(&dbEntry)

	return entry, nil
}

func (r *announcementRepository) GetAll(ctx context.Context) (domain.AnnouncementList, error) {
	var dbEntries []dbAnnouncement

	err := r.db.WithContext(ctx).
		Order("id ASC").
		Find(&dbEntries).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch announcements: %w", err)
	}

	return toDomainAnnouncementList(&dbEntries), nil
}

func (r *announcementRepository) applyFilters(db *gorm.DB, f repository.AnnouncementListFilters) *gorm.DB {
	if f.Query != nil {
		likeQuery := fmt.Sprintf("%%%s%%", *f.Query)
		db = db.Where("title LIKE ? OR body LIKE ?", likeQuery, likeQuery)
	}

	if f.ID != nil {
		db = db.Where("id = ?", *f.ID)
	}

	if f.Hidden != nil {
		db = db.Where("is_hidden = ?", *f.Hidden)
	}

	return db
}
