package models

type FilterConnection struct {
	ID       uint `form:"filter_connection_id" yaml:"filter_connection_id" json:"filter_connection_id"`
	WatchID  uint `form:"connection_watch_id" gorm:"index" yaml:"connection_watch_id" json:"connection_watch_id" binding:"required"`
	OutputID uint `form:"filter_output_id" gorm:"index" yaml:"filter_output_id" json:"filter_output_id" binding:"required"`
	InputID  uint `form:"filter_input_id" gorm:"index" yaml:"filter_input_id" json:"filter_input_id" binding:"required"`
}
