package users

import "github.com/gofrs/uuid"

// User represents a user entity.
type User struct {
	UUID uuid.UUID `gorm:"type:uuid;notNull"`
}
