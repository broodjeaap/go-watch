package models

import (
	"fmt"
	"html"

	"github.com/robfig/cron/v3"
)

type FilterID uint
type FilterType string
type Filter struct {
	ID        FilterID    `form:"filter_id" yaml:"filter_id" json:"filter_id"`
	WatchID   WatchID     `form:"filter_watch_id" gorm:"index" yaml:"filter_watch_id" json:"filter_watch_id" binding:"required"`
	Name      string      `form:"filter_name" gorm:"index" yaml:"filter_name" json:"filter_name" binding:"required" validate:"min=1"`
	X         int         `form:"x" yaml:"x" json:"x" validate:"default=0"`
	Y         int         `form:"y" yaml:"y" json:"y" validate:"default=0"`
	Type      FilterType  `form:"filter_type" yaml:"filter_type" json:"filter_type" binding:"required" validate:"oneof=url xpath json css replace match substring math store condition cron"`
	Var1      string      `form:"var1" yaml:"var1" json:"var1" binding:"required"`
	Var2      string      `form:"var2" yaml:"var2" json:"var2"`
	Parents   []*Filter   `gorm:"-:all"`
	Children  []*Filter   `gorm:"-:all"`
	Results   []string    `gorm:"-:all"`
	Logs      []string    `gorm:"-:all"`
	CronEntry *cron.Entry `gorm:"-:all"`
}

func (filter *Filter) Log(v ...any) {
	filter.Logs = append(filter.Logs, html.EscapeString(fmt.Sprint(v...)))
}
