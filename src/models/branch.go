package models

import (
	"time"

	"gorm.io/gorm"
)


type Branch struct {
	gorm.Model
	Name        string      `gorm:"type:varchar(80);"`
	Hash        string      `gorm:"type:varchar(80);"`
	CommittedAt *time.Time
	CommitMsg   *string     `gorm:"type:text;"`
	IsTag       bool
	RepositoryId	uint
	Repository	Repository
}
