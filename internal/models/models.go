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

type GeoEntity struct {
	Id       int             `json:"-"`
	Code     string          `json:"code"`
	Name     string          `json:"name"`
	Type     string          `json:"type"`
	Children []GeoEntity     `json:"children"`
	Metadata []GeoEntityMeta `json:"metadata"`
}

type GeoEntityMeta struct {
	GeoEntityId int    `json:"-"`
	Key         string `json:"key"`
	Value       string `json:"value"`
}

type GeoRelationship struct {
	ParentId int        `json:"parent_id"`
	ChildId  int        `json:"child_id"`
	Since    time.Time  `json:"since"`
	Until    *time.Time `json:"until,omitempty"`
}

type Indicator struct {
	Id          int      `json:"-"`
	Code        string   `json:"code"`
	Name        string   `json:"name"`
	Category    string   `json:"category"`
	Unit        string   `json:"unit"`
	Description string   `json:"description"`
	DataType    DataType `json:"data_type"`
}

type DataSource struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	URL         string    `json:"url"`
	Description string    `json:"description"`
	LastUpdated time.Time `json:"last_updated"`
}

type DataPoint struct {
	Id           int       `json:"-"`
	GeoEntityId  int       `json:"geo_entity_id"`
	IndicatorId  int       `json:"indicator_id"`
	SourceId     int       `json:"source_id"`
	Date         time.Time `json:"date"`
	NumericValue *float64  `json:"numeric_value,omitempty"`
	TextValue    *string   `json:"text_value,omitempty"`
}
