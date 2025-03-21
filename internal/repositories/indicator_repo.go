package repository

import (
	"core/internal/models"
	"database/sql"
	"errors"
)

type IndicatorRepository struct {
	DB *sql.DB
}

func NewIndicatorRepository(db *sql.DB) *IndicatorRepository {
	return &IndicatorRepository{DB: db}
}

func (r *IndicatorRepository) GetIndicators(limit int, offset int) ([]models.Indicator, error) {
	query := "SELECT * FROM indicators ORDER BY id LIMIT $1 OFFSET $2"
	rows, err := r.DB.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	indicators := []models.Indicator{}
	for rows.Next() {
		var ind models.Indicator
		if err := rows.Scan(&ind.Id, &ind.Name, &ind.Category, &ind.Unit, &ind.Description); err != nil {
			return nil, err
		}
		indicators = append(indicators, ind)
	}
	return indicators, nil

}

func (r *IndicatorRepository) GetIndicatorById(id int) (*models.Indicator, error) {
	query := "SELECT * FROM indicators WHERE id = $1 LIMIT 1"
	row := r.DB.QueryRow(query, id)

	var ind models.Indicator
	err := row.Scan(&ind.Id, &ind.Name, &ind.Category, &ind.Unit, &ind.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("indicator not found")
		}
		return nil, err
	}
	return &ind, nil
}
