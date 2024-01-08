package po

import (
	"time"
)

// Base is the base model for all data model.
type Base struct {
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp with time zone" json:"updated_at" sig:"-"`
	CreatedAt time.Time `gorm:"type:timestamp with time zone;default:CURRENT_TIMESTAMP" json:"created_at" sig:"-"`
}
