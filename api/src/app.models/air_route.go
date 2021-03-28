package app_models

import (
	base "app.base"
	"bytes"
	"context"
	"fmt"
	"github.com/jackc/pgx/pgxpool"
)

type AirRouteRequest struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

func (ar *AirRouteRequest) Load(data []byte) error {
	return base.JsonProcessor.NewDecoder(bytes.NewReader(data)).Decode(ar)
}

type AirRouteResponse struct {
	Source      *Location `json:"source"`
	Destination *Location `json:"destination"`
}

func (arr *AirRouteResponse) Marshal() ([]byte, error) {
	data, err := base.JsonProcessor.Marshal(arr)
	return data, err
}

func (arr *AirRouteResponse) GetLocations(ctx context.Context, req *AirRouteRequest) error {
	sqlSelectWhereSS := " iata=$1 "
	sqlSelectWhereDD := " iata=$2 "

	if len(req.Source) > 3 {
		sqlSelectWhereSS = " icao=$1 "
	}

	if len(req.Destination) > 3 {
		sqlSelectWhereDD = " icao=$2 "
	}
	sqlSelect := `
		WITH ss as (
			SELECT latitude as s_latitude, longitude as s_longitude
			FROM airports
			WHERE %s
			LIMIT 1
		), dd as (
			SELECT latitude as d_latitude, longitude as d_longitude
			FROM airports
			WHERE %s
			LIMIT 1
		)
		SELECT s_latitude::float8, s_longitude::float8, d_latitude::float8, d_longitude::float8
		FROM ss, dd
		LIMIT 1;`

	sqlSelect = fmt.Sprintf(sqlSelect, sqlSelectWhereSS, sqlSelectWhereDD)
	pool, ok := ctx.Value("pgxpool").(*pgxpool.Pool)
	if !ok {
		return fmt.Errorf("error get pgxpool from context")
	}

	source := Location{}
	destination := Location{}
	err := pool.QueryRow(ctx, sqlSelect, req.Source, req.Destination).Scan(
		&source.Latitude,
		&source.Longitude,
		&destination.Latitude,
		&destination.Longitude,
	)
	if err != nil {
		log := base.GetLog(ctx, nil)
		log.Error("AirRouteResponse.GetLocations error get api route waypoints: ", err)
		return fmt.Errorf("error get api route waypoints")
	}
	arr.Source = &source
	arr.Destination = &destination
	return nil
}
