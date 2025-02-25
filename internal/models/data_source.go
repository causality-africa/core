package models

type DataSource struct {
	Id          int    `json:"id"`
	Name        string `json:"string"`
	Url         string `json:"url"`
	Descriptoin string `json:"description"`
}
