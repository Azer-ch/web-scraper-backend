package types

import (
	"time"
	"gorm.io/gorm"
)

type Analysis struct {
	ID         uint   `gorm:"primaryKey"`
	URL        string `gorm:"size:2048"`
	URLHash    []byte `gorm:"uniqueIndex;size:32"`
	Result     string `gorm:"type:json"`
	AnalyzedAt time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}
