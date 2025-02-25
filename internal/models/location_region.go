package models

import "time"

type LocationRegion struct {
	LocationId int       `json:"location_id"`
	RegionId   int       `json:"region_id"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
}
