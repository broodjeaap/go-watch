package main

import (
	"gorm.io/gorm"
)

type Watch struct {
	gorm.Model
	Name     string `form:"watch_name" yaml:"watch_name" binding:"required" validate:"min=1"`
	Interval int    `form:"interval" yaml:"interval" binding:"required"`
	Filters  []Filter
}

type Filter struct {
	gorm.Model
	WatchID  uint `form:"filter_watch_id" yaml:"filter_watch_id" binding:"required"`
	Watch    Watch
	ParentID *uint    `form:"parent_id" yaml:"parent_id"`
	Parent   *Filter  `form:"parent_id" yaml:"parent_id"`
	Name     string   `form:"filter_name" yaml:"filter_name" binding:"required" validate:"min=1"`
	Type     string   `form:"filter_type" yaml:"filter_type" binding:"required" validate:"oneof=url xpath json css replace match substring"`
	Var1     string   `form:"var1" yaml:"var1" binding:"required"`
	Var2     *string  `form:"var2" yaml:"var2"`
	Var3     *string  `form:"var3" yaml:"var3"`
	Filters  []Filter `gorm:"-:all"`
	results  []string `gorm:"-:all"`
}
