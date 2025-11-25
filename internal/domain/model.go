package domain

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type User struct {
	ID              string         `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()" json:"id"`

	FullName        string         `gorm:"type:varchar(120);not null" json:"full_name"`

	Email           string         `gorm:"type:citext;uniqueIndex;not null" json:"email"` // Unique + indexed

	Phone           *string        `gorm:"type:varchar(32);index" json:"phone,omitempty"`  // Add index for search

	PasswordHash    *string        `gorm:"type:text" json:"-"`

	Role            string         `gorm:"type:varchar(32);default:user;index" json:"role"` // Index for role filtering

	Status          string         `gorm:"type:varchar(16);default:active;index" json:"status"` // Index for active users

	// Optional Composite Index: speeds up lookups by email + status
	// Example: active user lookup during login
	// gorm:"index:idx_email_status,priority:2"
	EmailStatus     string         `gorm:"-"` // Not stored, only illustrating composite index approach

	EmailVerifiedAt *time.Time     `json:"email_verified_at,omitempty"`
	LastLoginAt     *time.Time     `json:"last_login_at,omitempty"`

	Provider        *string        `gorm:"type:varchar(32)" json:"provider,omitempty"`

	Metadata        datatypes.JSON `gorm:"type:jsonb;default:'{}'" json:"metadata"`

	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`

	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // Soft-delete index
}
