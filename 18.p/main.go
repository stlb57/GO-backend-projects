package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Category struct {
	ID          int64
	Name        string
	Description string
	CreatedAt   time.Time
}

func createCategory(ctx context.Context, pool *pgxpool.Pool, name string, description string) (int64, error) {
	var id int64

	err := pool.QueryRow(
		ctx,
		`INSERT INTO categories (name, description)
		 VALUES ($1, $2)
		 RETURNING id`,
		name,
		description,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func getCategory(ctx context.Context, pool *pgxpool.Pool, id int64) (Category, error) {
	var category Category

	err := pool.QueryRow(
		ctx,
		`SELECT id, name, description, created_at
		 FROM categories
		 WHERE id = $1`,
		id,
	).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
	)

	if err != nil {
		return Category{}, err
	}

	return category, nil
}

func listCategories(ctx context.Context, pool *pgxpool.Pool) ([]Category, error) {
	rows, err := pool.Query(
		ctx,
		`SELECT id, name, description, created_at
		 FROM categories
		 ORDER BY id`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []Category

	for rows.Next() {
		var category Category

		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func updateCategory(ctx context.Context, pool *pgxpool.Pool, id int64, name string, description string) error {
	_, err := pool.Exec(
		ctx,
		`UPDATE categories
		 SET name = $1,
		     description = $2
		 WHERE id = $3`,
		name,
		description,
		id,
	)

	return err
}

func deleteCategory(ctx context.Context, pool *pgxpool.Pool, id int64) error {
	_, err := pool.Exec(
		ctx,
		`DELETE FROM categories
		 WHERE id = $1`,
		id,
	)

	return err
}

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
	if err != nil {
		fmt.Println("Couldnt read db")
		os.Exit(1)
	}
	str := string(arr)
	_, err = pool.Exec(ctx, str)
	if err != nil {
		fmt.Println("Couldnt exec db")
		os.Exit(1)
	}

}
