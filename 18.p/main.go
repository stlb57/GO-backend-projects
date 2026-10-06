package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	pool, err := pgxpool.New(
		ctx,
		"postgres://postgres:admin@localhost:5432/inventory_db",
	)
	if err != nil {
		fmt.Println("connection failed:", err)
		os.Exit(1)
	}
	defer pool.Close()

	err = pool.Ping(ctx)
	if err != nil {
		fmt.Println("ping failed:", err)
		os.Exit(1)
	}

	arr, err := os.ReadFile("migrations/001_init.sql")
	str := string(arr)
	_, err = pool.Exec(ctx, str)

}
