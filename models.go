package main

import (
	"gorm.io/gorm"
)

type Watch struct {
	gorm.Model
	Name     string `form:"watch_name" yaml:"watch_name" binding:"required" validate:"min=1"`
	Interval int    `form:"interval" yaml:"interval" binding:"required"`
	URLs     []URL
}

type URL struct {
	gorm.Model
	WatchID      uint   `form:"url_watch_id" yaml:"url_watch_id" binding:"required"`
	Watch        *Watch `form:"watch" yaml:"watch" validate:"omitempty"`
	Name         string `form:"url_name" yaml:"url_name" binding:"required" validate:"min=1"`
	URL          string `form:"url" yaml:"url" binding:"required,url" validate:"min=1"`
	GroupFilters []FilterGroup
}

type FilterGroup struct {
	gorm.Model
	URLID   uint `form:"group_url_id" yaml:"group_url_id" binding:"required"`
	URL     *URL
	Name    string `form:"group_name" yaml:"group_name" binding:"required" validate:"min=1"`
	Type    string `form:"group_type" yaml:"group_type" binding:"required" validate:"oneof=diff enum number bool"`
	Filters []Filter
}

type Filter struct {
	gorm.Model
	FilterGroupID uint `form:"filter_group_id" yaml:"filter_group_id" binding:"required"`
	FilterGroup   *FilterGroup
	Name          string `form:"filter_name" yaml:"filter_name" binding:"required" validate:"min=1"`
	Type          string `form:"filter_type" yaml:"filter_type" binding:"required" validate:"oneof=xpath css replace match substring"`
	From          string `form:"from" yaml:"from" binding:"required"`
	To            string `form:"to" yaml:"to" binding:"required"`
}
