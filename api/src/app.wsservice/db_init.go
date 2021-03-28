package app_wsservice

import (
	"context"
	"fmt"
	"github.com/jackc/pgx"
)

func testTableExist(ctx context.Context, c *pgx.Conn, tablename string) (bool, error) {
	sql := `SELECT to_regclass($1);`
	var name *string
	err := c.QueryRow(ctx, sql, tablename).Scan(&name)
	switch err {
	case nil:
		if name == nil {
			return false, nil
		}
		return *name == tablename, nil
	case pgx.ErrNoRows:
		return false, nil
	default:
		return false, fmt.Errorf("testTableExist error answer: %s", err)
	}
}

func initAirportsTable(ctx context.Context, c *pgx.Conn) error {
	ok, err := testTableExist(ctx, c, "tbfm_aircraft_data")
	if err != nil {
		fmt.Println("initAirportsTable error: ", err)
		return err
	}
	if ok {
		return nil
	}

	sqlCreate := `
		CREATE TABLE IF NOT EXISTS airports (
			id			bigserial,
			name        text,
			city        text,
			country     text,
			iata        text,
			icao        text,
			Latitude    float,
			Longitude   float,
			Altitude    float,
			Timezone    int,
			DST         text,
			Tz          text,
			Type        text,
			Source      text
		);`
	sqlCreateIndexIata := `
		CREATE INDEX iata_airports_idx ON airports (iata);`

	sqlCreateIndexIcao := `
		CREATE INDEX icao_airports_idx ON airports (icao);`

	tx, err := c.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("initAirportsTable error start tx")
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, sqlCreate)
	if err != nil {
		return fmt.Errorf("initAirportsTable error create table: %s", err)
	}

	_, err = tx.Exec(ctx, sqlCreateIndexIata)
	if err != nil {
		return fmt.Errorf("initAirportsTable error create iata index: %s", err)
	}

	_, err = tx.Exec(ctx, sqlCreateIndexIcao)
	if err != nil {
		return fmt.Errorf("initAirportsTable error create icao index: %s", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("initAirportsTable error commit: %s", err)
	}

	return nil
}
