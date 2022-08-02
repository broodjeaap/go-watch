package main

// todo move bind/validate funcs here also

type FilterGroupUpdate struct {
	ID uint `form:"group_id" binding:"required"`
	// TODO add group_id for ezpz
	Name string `form:"group_name" binding:"required" validate:"min=1"`
	Type string `form:"group_type" binding:"required" validate:"oneof=diff enum number"`
}

type FilterUpdate struct {
	ID      uint   `form:"filter_id" binding:"required"`
	GroupID uint   `form:"filter_group_id" binding:"required"`
	Name    string `form:"name" binding:"required" validate:"min=1"`
	Type    string `form:"filter_type" binding:"required"`
	From    string `form:"from" binding:"required"`
	To      string `form:"to" binding:"required"`
}
