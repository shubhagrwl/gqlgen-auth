package dbmodels

import (
	"time"

	"github.com/google/uuid"
)

// Different tables that are used by the app.
const (
	TABLEUSERDETAIL = "user_detail"

	TABLEUSERDETAIL_COLUMN__ID              = "id"
	TABLEUSERDETAIL_COLUMN__Password        = "password"
	TABLEUSERDETAIL_COLUMN__FirstName       = "first_name"
	TABLEUSERDETAIL_COLUMN__LastName        = "last_name"
	TABLEUSERDETAIL_COLUMN__Email           = "email"
	TABLEUSERDETAIL_COLUMN__IsEmailVerified = "is_email_verified"
	TABLEUSERDETAIL_COLUMN__Dob             = "dob"
	TABLEUSERDETAIL_COLUMN__LastLoginAt     = "last_login_at"
	TABLEUSERDETAIL_COLUMN__Active          = "active"
	TABLEUSERDETAIL_COLUMN__CreatedAt       = "created_at"
	TABLEUSERDETAIL_COLUMN__UpdatedAt       = "updated_at"
)

// UserDetail Table structure
type UserDetail struct {
	ID              *uuid.UUID `gorm:"primary_key"`
	Password        *string    `gorm:"Column:password"`
	FirstName       *string    `gorm:"Column:first_name"`
	LastName        *string    `gorm:"Column:last_name"`
	Email           *string    `gorm:"Column:email"`
	IsEmailVerified *bool      `gorm:"Column:is_email_verified"`
	Dob             *time.Time `gorm:"Column:dob"`
	LastLoginAt     *string    `gorm:"Column:last_login_at"`
	Active          *bool      `gorm:"Column:active"`
	CreatedAt       *string    `gorm:"Column:created_at"`
	UpdatedAt       *string    `gorm:"Column:updated_at"`
}
