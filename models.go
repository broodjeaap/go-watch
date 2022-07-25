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
	WatchID uint   `form:"url_watch_id" yaml:"url_watch_id" binding:"required"`
	Watch   *Watch `form:"watch" yaml:"watch" validate:"omitempty"`
	Name    string `form:"url_name" yaml:"url_name" binding:"required" validate:"min=1"`
	URL     string `form:"url" yaml:"url" binding:"required,url" validate:"min=1"`
	Queries []Query
}

type Query struct {
	gorm.Model
	URLID   uint `form:"query_url_id" yaml:"query_url_id" binding:"required"`
	URL     *URL
	Name    string `form:"query_name" yaml:"query_name" binding:"required" validate:"min=1"`
	Type    string `form:"query_type" yaml:"query_type" binding:"required" validate:"oneof=css xpath regex json"`
	Query   string `form:"query" yaml:"query" binding:"required"`
	Filters []Filter
}

type Filter struct {
	gorm.Model
	QueryID uint `form:"filter_query_id" yaml:"filter_query_id" binding:"required"`
	Query   *Query
	Name    string `form:"filter_name" yaml:"filter_name" binding:"required" validate:"min=1"`
	Type    string `form:"filter_type" yaml:"filter_type" binding:"required" validate:"oneof=replace regex substring"`
	From    string `form:"from" yaml:"from" binding:"required"`
	To      string `form:"to" yaml:"to" binding:"required"`
}
