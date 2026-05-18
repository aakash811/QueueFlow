package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func ConnectDatabase(databaseURL string) error {
	pool, err := pgxpool.New(context.Background(), databaseURL)

	if err != nil {
		return err

	}

	err = pool.Ping(context.Background())

	if err != nil {
		return err
	}

	fmt.Println("database connected successfully")
	
	DB = pool
	return nil
}