package models

import (
	"encoding/json"
	"fmt"
	"time"
)

// DataType represents the type of data an indicator can have
type DataType string

const (
	NumericDataType     DataType = "numeric"
	CategoricalDataType DataType = "categorical"
	BooleanDataType     DataType = "boolean"
)

func (dt DataType) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(dt))
}

func (dt *DataType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch s {
	case string(NumericDataType), string(CategoricalDataType), string(BooleanDataType):
		*dt = DataType(s)
		return nil
	default:
		return fmt.Errorf("invalid indicator data type: %s", s)
	}
}

// EntityType can be either location or region
type EntityType string

const (
	LocationEntityType EntityType = "location"
	RegionEntityType   EntityType = "region"
)

func (et EntityType) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(et))
}

func (et *EntityType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch s {
	case string(LocationEntityType), string(RegionEntityType):
		*et = EntityType(s)
		return nil
	default:
		return fmt.Errorf("invalid entity type: %s", s)
	}
}

type Location struct {
	Id         int     `json:"id"`
	Name       string  `json:"name"`
	Code       string  `json:"code"`
	AdminLevel int     `json:"admin_level"`
	ParentID   *int    `json:"parent_id,omitempty"`
	Map        *string `json:"map,omitempty"`
}

type Indicator struct {
	Id          int      `json:"id"`
	Name        string   `json:"name"`
	Code        string   `json:"code"`
	Category    string   `json:"category"`
	Unit        *string  `json:"unit,omitempty"`
	Description *string  `json:"description,omitempty"`
	DataType    DataType `json:"data_type"`
}

type DataSource struct {
	Id          int        `json:"id"`
	Name        string     `json:"name"`
	URL         *string    `json:"url,omitempty"`
	Description *string    `json:"description,omitempty"`
	Date        *time.Time `json:"date,omitempty"`
}

type DataPoint struct {
	Id           int        `json:"id"`
	EntityType   EntityType `json:"entity_type"`
	EntityID     int        `json:"entity_id"`
	IndicatorId  int        `json:"indicator_id"`
	SourceId     int        `json:"source_id"`
	Date         time.Time  `json:"date"`
	NumericValue *float64   `json:"numeric_value,omitempty"`
	TextValue    *string    `json:"text_value,omitempty"`
}

type Region struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Code        string  `json:"code"`
	Description *string `json:"description,omitempty"`
}

type LocationInRegion struct {
	LocationId int        `json:"location_id"`
	RegionId   int        `json:"region_id"`
	JoinDate   time.Time  `json:"join_date"`
	ExitDate   *time.Time `json:"exit_date,omitempty"`
}
