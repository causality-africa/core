package models

type Location struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	AdminLevel int    `json:"admin_level"`
	ParentId   int    `json:"parent_id"`
	IsoCode    string `json:"iso_code"`
	Map        string `json:"map"`
}
