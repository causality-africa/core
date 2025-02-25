package models

type Indicator struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Unit        string `json:"unit"`
	Description string `json:"description"`
}
