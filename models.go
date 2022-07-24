package main

import (
	"gorm.io/gorm"
)

type Watch struct {
	gorm.Model
	Name     string `form:"name" yaml:"name" binding:"required" validate:"min=1"`
	Interval int    `form:"interval" yaml:"interval" binding:"required"`
	URLs     []URL
}

type URL struct {
	gorm.Model
	WatchID uint   `form:"watch_id" yaml:"watch_id" binding:"required"`
	Watch   *Watch `form:"watch" yaml:"watch" validate:"omitempty"`
	Name    string `form:"name" yaml:"name" binding:"required"`
	URL     string `form:"url" yaml:"url" binding:"required,url"`
	Queries []Query
}

type Query struct {
	gorm.Model
	URLID   uint `form:"url_id" yaml:"url_id" binding:"required"`
	URL     URL
	Name    string `form:"name" yaml:"name" binding:"required"`
	Type    string `form:"type" yaml:"type" binding:"required"`
	Query   string `form:"query" yaml:"query" binding:"required"`
	Filters []Filter
}

type Filter struct {
	gorm.Model
	QueryID uint `form:"query_id" yaml:"query_id" binding:"required"`
	Query   Query
	Name    string `form:"name" yaml:"name" binding:"required"`
	Type    string `form:"type" yaml:"type" binding:"required"`
	From    string `form:"from" yaml:"from" binding:"required"`
	To      string `form:"to" yaml:"to" binding:"required"`
}
