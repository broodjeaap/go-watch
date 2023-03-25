package models

import "time"

type FilterOutputID uint
type FilterOutput struct {
	ID      FilterOutputID `yaml:"filter_output_id" json:"filter_output_id"`
	WatchID WatchID        `yaml:"filter_output_watch_id" gorm:"index" json:"filter_output_watch_id"`
	Name    string         `yaml:"filter_output_name" json:"filter_output_name"`
	Value   string         `yaml:"filter_output_value" json:"filter_output_value"`
	Time    time.Time      `yaml:"filter_output_time" json:"filter_output_time"`
}
