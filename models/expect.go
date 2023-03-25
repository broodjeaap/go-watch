package models

import "time"

type ExcpectFailID uint
type ExpectFail struct {
	ID      ExcpectFailID `yaml:"expect_fail_id" json:"expect_fail_id"`
	WatchID WatchID       `yaml:"expect_fail_watch_id" gorm:"index" json:"expect_fail_watch_id"`
	Name    string        `yaml:"expect_fail_name" json:"expect_fail_name"`
	Time    time.Time     `yaml:"expect_fail_time" json:"expect_fail_time"`
}
