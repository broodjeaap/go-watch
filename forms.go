package main

type QueryUpdate struct {
	ID    uint   `form:"query_id" binding:"required"`
	Name  string `form:"query_name" binding:"required" validate:"min=1"`
	Type  string `form:"query_type" binding:"required" validate:"oneof=css xpath regex json"`
	Query string `form:"query" binding:"required"`
}
