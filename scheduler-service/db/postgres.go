package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectDatabase(databaseURL string) error {

	dbpool, err := pgxpool.New(
		context.Background(),
		databaseURL,
	)

	if err != nil {
		return err
	}

	err = dbpool.Ping(context.Background())

	if err != nil {
		return err
	}

	fmt.Println("scheduler database connected")

	DB = dbpool

	return nil
}
