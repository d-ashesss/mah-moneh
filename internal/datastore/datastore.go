package datastore

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"time"
)

// Model defines fields common for most models.
type Model struct {
	UUID      uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
