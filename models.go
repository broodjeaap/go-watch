package main

import (
	"gorm.io/gorm"
)

type Watch struct {
	gorm.Model
	Name     string
	Interval int
	URLs     []URL
}

type URL struct {
	gorm.Model
	WatchID uint
	Watch   Watch
	Name    string
	URL     string
	Queries []Query
}

type Query struct {
	gorm.Model
	URLID   uint
	URL     URL
	Name    string
	Type    string
	Query   string
	Filters []Filter
}

type Filter struct {
	gorm.Model
	QueryID uint
	Query   Query
	Name    string
	Type    string
	From    string
	To      string
}
