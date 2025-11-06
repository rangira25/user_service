package domain

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type User struct {
	ID              string         `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`
	FullName        string         `gorm:"type:varchar(120);not null" json:"full_name"`
	Email           string         `gorm:"type:citext;uniqueIndex;not null" json:"email"`
	Phone           *string        `gorm:"type:varchar(32)" json:"phone,omitempty"`
	PasswordHash    *string        `gorm:"type:text" json:"-"`
	Role            string         `gorm:"type:varchar(32);default:user" json:"role"`
	Status          string         `gorm:"type:varchar(16);default:active" json:"status"`
	EmailVerifiedAt *time.Time     `json:"email_verified_at,omitempty"`
	LastLoginAt     *time.Time     `json:"last_login_at,omitempty"`
	Provider        *string        `gorm:"type:varchar(32)" json:"provider,omitempty"`
	Metadata        datatypes.JSON `gorm:"type:jsonb;default:'{}'" json:"metadata"` 
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
