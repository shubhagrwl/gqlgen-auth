package onboarduser

import "time"

// Different tables that are used by the app.
const (
	TABLE_ONBOARD_USER = "onboard_user"

	TABLE_ONBOARD_USER_COLUMN__ID         = "id"
	TABLE_ONBOARD_USER_COLUMN__Email      = "email"
	TABLE_ONBOARD_USER_COLUMN__Verified   = "verified"
	TABLE_ONBOARD_USER_COLUMN__CreatedAt  = "created_at"
	TABLE_ONBOARD_USER_COLUMN__VerifiedAt = "verified_at"
)

// OnboardUser Table structure
type OnboardUser struct {
	ID        *int       `gorm:"id"`
	Email     *string    `gorm:"Column:email"`
	Verified  *bool      `gorm:"Column:verified"`
	CreatedAt *time.Time `gorm:"Column:created_at"`
	UpdatedAt *time.Time `gorm:"Column:verified_at"`
}
