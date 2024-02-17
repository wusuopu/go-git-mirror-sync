package models

import (
	"time"

	"gorm.io/gorm"
)


type Mirror struct {
	gorm.Model
	Name        string      `gorm:"type:varchar(80);"`
	Alias       string      `gorm:"type:varchar(80);"`
	Url         string      `gorm:"type:varchar(255);"`
	AuthType    string      `gorm:"type:varchar(15);"`      // password | sshkey
	Username    *string     `gorm:"type:varchar(80);"`
	Password    *string     `gorm:"type:varchar(80);"`
	SSHKey      *string     `gorm:"type:varchar(80);"`
	PushedAt    *time.Time
	LastError   *string     `gorm:"type:text;"`
	RepositoryId	uint
	Repository	Repository
}
