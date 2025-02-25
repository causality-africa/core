package models

import "time"

type DataPoint struct {
	Id          int       `json:"id"`
	LocationId  int       `json:"location_id"`
	IndicatorId int       `json:"indicator_id"`
	SourceId    int       `json:"source_id"`
	Date        time.Time `json:"date"`
	Value       float64   `json:"value"`
	IsEstimated bool      `json:"is_estimated"`
}
