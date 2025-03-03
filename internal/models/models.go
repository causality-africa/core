package models

import (
	"database/sql"
	"time"
)

type DataPoint struct {
	Id          int       `json:"id"`
	LocationId  int       `json:"location_id"`
	IndicatorId int       `json:"indicator_id"`
	SourceId    int       `json:"source_id"`
	Date        time.Time `json:"date"`
	Value       float64   `json:"value"`
	IsEstimated bool      `json:"is_estimated"`
}

type DataSource struct {
	Id          int    `json:"id"`
	Name        string `json:"string"`
	Url         string `json:"url"`
	Descriptoin string `json:"description"`
}

type Indicator struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Unit        string `json:"unit"`
	Description string `json:"description"`
}

type LocationRegion struct {
	LocationId int       `json:"location_id"`
	RegionId   int       `json:"region_id"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
}

type Location struct {
	Id         int           `json:"id"`
	Name       string        `json:"name"`
	AdminLevel int           `json:"admin_level"`
	ParentId   sql.NullInt32 `json:"parent_id"`
	Code       string        `json:"code"`
	Map        string        `json:"map"`
}

type Region struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
