package models

import "time"

type ExpectFail struct {
	ID      uint      `yaml:"expect_fail_id" json:"expect_fail_id"`
	WatchID uint      `yaml:"expect_fail_watch_id" gorm:"index" json:"expect_fail_watch_id"`
	Name    string    `yaml:"expect_fail_name" json:"expect_fail_name"`
	Time    time.Time `yaml:"expect_fail_time" json:"expect_fail_time"`
}
