package app_models

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/pgxpool"
)

type Location struct {
	City      string  `json:"city"`
	Country   string  `json:"country"`
	Iata      string  `json:"iata"`
	Icao      string  `json:"icao"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
	Altitude  int     `json:"alt"`
	Timezone  int     `json:"tz"`
}

func (arr *Location) GetLocationGeo(ctx context.Context, code string) error {
	sqlSelectWhere := " iata=$1 "
	sqlSelect := `
		SELECT latitude, latitude 
		FROM airports
		WHERE %
		LIMIT 1;`
	if len(code) > 3 {
		sqlSelectWhere = " icao=$1 "
	}
	sqlSelect = fmt.Sprintf(sqlSelect, sqlSelectWhere)

	pool, ok := ctx.Value("pgxpool").(*pgxpool.Pool)
	if !ok {
		return fmt.Errorf("error get pgxpool from context")
	}

	return pool.QueryRow(ctx, sqlSelect, code).Scan(&arr.Latitude, &arr.Longitude)
}
