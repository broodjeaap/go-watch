package models

import (
	"github.com/robfig/cron/v3"
)

type Watch struct {
	ID        uint        `form:"watch_id" yaml:"watch_id"`
	Name      string      `form:"watch_name" gorm:"index" yaml:"watch_name" binding:"required" validate:"min=1"`
	CronEntry *cron.Entry `gorm:"-:all"`
	LastValue string      `gorm:"-:all"`
}
