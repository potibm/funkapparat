package gorm

import (
	"context"
	"testing"

	"github.com/potibm/funkapparat/internal/app/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type TestLocation struct {
	ID         uint `gorm:"primaryKey"`
	Name       string
	CreatedBy  string
	ModifiedBy string
	DeletedBy  *string
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func TestAuditCallbacks(t *testing.T) {
	db, err := NewSqliteInMemoryStore()
	require.NoError(t, err)
	defer db.Close()

	gormDB := db.db

	err = RegisterAuditCallbacks(gormDB)
	require.NoError(t, err)

	err = gormDB.AutoMigrate(&TestLocation{})
	require.NoError(t, err)

	creatorID := "user-creator-123"
	updaterID := "user-updater-456"
	deleterID := "user-deleter-789"

	t.Run("BeforeCreate sets CreatedBy and ModifiedBy", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), domain.UserIDKey, creatorID)

		loc := TestLocation{Name: "Main Stage"}
		err := gormDB.WithContext(ctx).Create(&loc).Error

		require.NoError(t, err)
		assert.Equal(t, creatorID, loc.CreatedBy, "CreatedBy should be set to the creator's user ID")
		assert.Equal(t, creatorID, loc.ModifiedBy, "ModifiedBy should be equal to CreatedBy initially")
	})

	t.Run("BeforeUpdate changes ONLY ModifiedBy", func(t *testing.T) {
		var loc TestLocation
		gormDB.First(&loc, 1)

		ctx := context.WithValue(context.Background(), domain.UserIDKey, updaterID)
		err := gormDB.WithContext(ctx).Model(&loc).Where("id = ?", loc.ID).Update("Name", "Updated Stage").Error

		require.NoError(t, err)
		assert.Equal(t, creatorID, loc.CreatedBy, "CreatedBy should NOT be changed")

		var updatedLoc TestLocation
		gormDB.First(&updatedLoc, 1)

		assert.Equal(t, creatorID, updatedLoc.CreatedBy, "CreatedBy should NOT be changed")
		assert.Equal(t, updaterID, updatedLoc.ModifiedBy, "ModifiedBy should be updated to the updater's user ID")
	})

	t.Run("BeforeDelete changes DeletedBy", func(t *testing.T) {
		var loc TestLocation
		err := gormDB.First(&loc, 1).Error
		require.NoError(t, err)

		ctx := context.WithValue(context.Background(), domain.UserIDKey, deleterID)
		err = gormDB.WithContext(ctx).Delete(&loc).Error

		require.NoError(t, err)

		var deletedLoc TestLocation
		err = gormDB.Unscoped().First(&deletedLoc, 1).Error

		require.NoError(t, err)
		assert.True(t, deletedLoc.DeletedAt.Valid, "GORM should have set DeletedAt for soft-deleted record")

		if deletedLoc.DeletedBy != nil {
			assert.Equal(t, deleterID, *deletedLoc.DeletedBy, "DeletedBy should be set to the deleter's user ID")
		} else {
			t.Errorf("DeletedBy should be set")
		}

		var untouchedReloaded TestLocation
		err = gormDB.First(&untouchedReloaded, 1).Error
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound, "Record should not be found without Unscoped()")
	})
}
