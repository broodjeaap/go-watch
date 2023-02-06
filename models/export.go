package models

type WatchExport struct {
	Filters     []Filter           `json:"filters"`
	Connections []FilterConnection `json:"connections"`
}
