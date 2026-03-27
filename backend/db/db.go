package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

var Conn *pgx.Conn

func InitDB() error {
	var err error
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return fmt.Errorf("DATABASE_URL не задан в .env")
	}

	Conn, err = pgx.Connect(context.Background(), dbURL)
	if err != nil {
		return fmt.Errorf("не удалось подключиться к БД: %w", err)
	}

	err = createTables()
	if err != nil {
		return fmt.Errorf("ошибка создания таблиц: %w", err)
	}

	fmt.Println("БД подключена успешно!")
	return nil
}

func createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS candidates (
		id          SERIAL PRIMARY KEY,
		name        TEXT NOT NULL,
		age         INT,
		essay       TEXT,
		experience  TEXT,
		motivation  TEXT,
		score       INT DEFAULT 0,
		explanation TEXT DEFAULT '',
		ai_detected BOOLEAN DEFAULT FALSE,
		status      TEXT DEFAULT 'pending',
		created_at  TIMESTAMP DEFAULT NOW()
	);`

	_, err := Conn.Exec(context.Background(), query)
	return err
}