package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type DB struct {
	pool *pgxpool.Pool
}

func NewDB(ctx context.Context) (*DB, error) {
	if err := godotenv.Load(); err != nil {
		   return nil, fmt.Errorf("error cargando .env: %w", err)
	}
	url := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("SSL_MODE"),
	)
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("error conect DB: %w", err)
	}
	return &DB{pool: pool}, nil
}

func (db *DB) InsertItem(item Item) error {
	targetFromStr := strings.ReplaceAll(item.TargetFrom, ",", "")
	targetFromStr = strings.Replace(targetFromStr, "$", "", 1)
	targetToStr := strings.ReplaceAll(item.TargetTo, ",", "")
	targetToStr = strings.Replace(targetToStr, "$", "", 1)
	targetFrom, err := strconv.ParseFloat(targetFromStr, 64)
	if err != nil {
		return fmt.Errorf("error parsing TargetFrom: %w", err)
	}
	targetTo, err := strconv.ParseFloat(targetToStr, 64)
	if err != nil {
		return fmt.Errorf("error parsing TargetTo: %w", err)
	}
	query := `
		INSERT INTO stock (
			ticker, target_from, target_to, company, action,
			brokerage, rating_from, rating_to, time
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);
	`
	_, err = db.pool.Exec(
		context.Background(),
		query,
		item.Ticker, targetFrom, targetTo,
		item.Company, item.Action, item.Brokerage,
		item.RatingFrom, item.RatingTo, item.Time,
	)
	return err
}
