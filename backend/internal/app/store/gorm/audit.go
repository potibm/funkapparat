package gorm

import (
	"github.com/potibm/funkapparat/internal/app/domain"
	"gorm.io/gorm"
	"reflect"
)

type AuditModel struct {
	CreatedBy  string  `json:"created_by"           gorm:"column:created_by;type:varchar(255);<-:create"`
	ModifiedBy string  `json:"modified_by"          gorm:"column:modified_by;type:varchar(255)"`
	DeletedBy  *string `json:"deleted_by,omitempty" gorm:"column:deleted_by;type:varchar(255);index"`
}

func RegisterAuditCallbacks(db *gorm.DB) error {
	cbCreate := db.Callback().Create()
	if cbCreate.Get("audit:before_create") == nil {
		if err := cbCreate.Before("gorm:create").Register("audit:before_create", beforeCreateCallback); err != nil {
			return err
		}
	}

	cbUpdate := db.Callback().Update()
	if cbUpdate.Get("audit:before_update") == nil {
		if err := cbUpdate.Before("gorm:update").Register("audit:before_update", beforeUpdateCallback); err != nil {
			return err
		}
	}

	cbDelete := db.Callback().Delete()
	if cbDelete.Get("audit:before_delete") == nil {
		if err := cbDelete.Before("gorm:delete").Register("audit:before_delete", beforeDeleteCallback); err != nil {
			return err
		}
	}

	return nil
}

func getUserIDFromContext(tx *gorm.DB) string {
	if tx.Statement == nil || tx.Statement.Context == nil {
		return ""
	}
	userID, ok := tx.Statement.Context.Value(domain.UserIDKey).(string)
	if !ok {
		return ""
	}
	return userID
}

func setAuditColumn(tx *gorm.DB, fieldName, value string) {
	if tx.Statement.Schema == nil {
		return
	}
	field := tx.Statement.Schema.LookUpField(fieldName)
	if field != nil {
		// Set in memory if applicable
		if tx.Statement.ReflectValue.IsValid() && tx.Statement.ReflectValue.CanSet() {
			_ = field.Set(tx.Statement.Context, tx.Statement.ReflectValue, value)
		}
		// Apply to SQL
		tx.Statement.SetColumn(field.Name, value, true)
	}
}

func beforeCreateCallback(tx *gorm.DB) {
	if userID := getUserIDFromContext(tx); userID != "" {
		setAuditColumn(tx, "CreatedBy", userID)
		setAuditColumn(tx, "ModifiedBy", userID)
	}
}

func beforeUpdateCallback(tx *gorm.DB) {
	if userID := getUserIDFromContext(tx); userID != "" {
		setAuditColumn(tx, "ModifiedBy", userID)
	}
}

func beforeDeleteCallback(tx *gorm.DB) {
	userID := getUserIDFromContext(tx)
	if userID == "" {
		return
	}

	if tx.Statement.Schema == nil {
		return
	}

	_, hasDeletedAt := tx.Statement.Schema.FieldsByDBName["deleted_at"]
	if !hasDeletedAt {
		return
	}

	// Because Soft Delete triggers an update implicitly *after* deleting we cannot easily add a clause.
	// Instead, just do a direct UPDATE query for the deleted_by before actual deletion proceeds.
	if tx.Statement.Dest != nil {
		var id interface{}

		val := reflect.ValueOf(tx.Statement.Dest)
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		if val.Kind() == reflect.Struct {
			idField := val.FieldByName("ID")
			if idField.IsValid() {
				id = idField.Interface()
			}
		} else if val.Kind() == reflect.Map {
			if mapID, ok := tx.Statement.Dest.(map[string]interface{})["id"]; ok {
				id = mapID
			}
		}

		if id != nil {
			tx.Statement.DB.Session(&gorm.Session{NewDB: true}).Table(tx.Statement.Table).Where("id = ?", id).UpdateColumn("deleted_by", userID)
		}
	}
}
