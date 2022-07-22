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
	Name    string
	URL     string
	Queries []Query
}

type Query struct {
	gorm.Model
	URLID uint
	Name  string
	Type  string
	Query string
}
