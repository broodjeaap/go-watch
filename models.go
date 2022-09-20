package main

import (
	"gorm.io/gorm"
)

type Watch struct {
	gorm.Model
	Name     string `form:"watch_name" yaml:"watch_name" binding:"required" validate:"min=1"`
	Interval int    `form:"interval" yaml:"interval" binding:"required"`
}

type Filter struct {
	gorm.Model
	WatchID uint     `form:"filter_watch_id" yaml:"filter_watch_id" binding:"required"`
	Name    string   `form:"filter_name" yaml:"filter_name" binding:"required" validate:"min=1"`
	X       int      `form:"x" yaml:"x" validate:"default=0"`
	Y       int      `form:"y" yaml:"y" validate:"default=0"`
	Type    string   `form:"filter_type" yaml:"filter_type" binding:"required" validate:"oneof=url xpath json css replace match substring"`
	Var1    string   `form:"var1" yaml:"var1" binding:"required"`
	Var2    string   `form:"var2" yaml:"var2"`
	Var3    string   `form:"var3" yaml:"var3"`
	Results []string `gorm:"-:all"`
}

type FilterConnection struct {
	gorm.Model
	WatchID  uint `form:"connection_watch_id" yaml:"connection_watch_id" binding:"required"`
	OutputID uint `form:"filter_output_id" yaml:"filter_output_id" binding:"required"`
	InputID  uint `form:"filter_output_id" yaml:"filter_output_id" binding:"required"`
}
