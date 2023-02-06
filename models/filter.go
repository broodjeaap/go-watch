package models

import (
	"fmt"
	"html"
)

type Filter struct {
	ID       uint      `form:"filter_id" yaml:"filter_id" json:"filter_id"`
	WatchID  uint      `form:"filter_watch_id" gorm:"index" yaml:"filter_watch_id" json:"filter_watch_id" binding:"required"`
	Name     string    `form:"filter_name" gorm:"index" yaml:"filter_name" json:"filter_name" binding:"required" validate:"min=1"`
	X        int       `form:"x" yaml:"x" json:"x" validate:"default=0"`
	Y        int       `form:"y" yaml:"y" json:"y" validate:"default=0"`
	Type     string    `form:"filter_type" yaml:"filter_type" json:"filter_type" binding:"required" validate:"oneof=url xpath json css replace match substring math store condition cron"`
	Var1     string    `form:"var1" yaml:"var1" json:"var1" binding:"required"`
	Var2     *string   `form:"var2" yaml:"var2" json:"var2"`
	Var3     *string   `form:"var3" yaml:"var3" json:"var3"`
	Parents  []*Filter `gorm:"-:all"`
	Children []*Filter `gorm:"-:all"`
	Results  []string  `gorm:"-:all"`
	Logs     []string  `gorm:"-:all"`
}

func (filter *Filter) Log(v ...any) {
	filter.Logs = append(filter.Logs, html.EscapeString(fmt.Sprint(v...)))
}
