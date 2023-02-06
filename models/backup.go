package models

type Backup struct {
	Watches     []Watch            `json:"watches"`
	Filters     []Filter           `json:"filters"`
	Connections []FilterConnection `json:"connections"`
	Values      []FilterOutput     `json:"values"`
}
