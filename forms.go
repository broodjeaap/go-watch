package main

type FilterGroupUpdate struct {
	ID   uint   `form:"group_id" binding:"required"`
	Name string `form:"group_name" binding:"required" validate:"min=1"`
	Type string `form:"group_type" binding:"required" validate:"oneof=diff enum number"`
}
