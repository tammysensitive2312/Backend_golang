package entities

import "time"

type Project struct {
	ID                int        `gorm:"primaryKey,autoIncrement"`
	Name              string     `gorm:"notnull"`
	Category          string     `gorm:"notnull"`
	ProjectSpend      int        `gorm:"notnull, default=0"`
	ProjectVariance   int        `gorm:"notnull, default=0"`
	RevenueRecognised int        `gorm:"notnull, default=0"`
	ProjectStartedAt  time.Time  `gorm:"notnull"`
	ProjectEndedAt    *time.Time `gorm:"type:datetime"`
	CreatedAt         time.Time  `gorm:"autoCreateTime"`
	UpdatedAt         time.Time  `gorm:"autoUpdateTime"`
	DeletedAt         *time.Time `gorm:"index"`
}
