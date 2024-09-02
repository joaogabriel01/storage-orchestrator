package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type PostgresStorageUnit struct {
	db *sql.DB
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "mysecretpassword"
	dbname   = "postgres"
)

func NewPostgresStorageUnit() *PostgresStorageUnit {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected!")

	return &PostgresStorageUnit{db: db}
}

func (p *PostgresStorageUnit) Save(ctx context.Context, _, userJson string) error {
	var user User
	err := json.Unmarshal([]byte(userJson), &user)
	if err != nil {
		return err
	}
	_, err = p.db.ExecContext(ctx, "INSERT INTO users (id, details) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET details = EXCLUDED.details", user.ID, user.Details)
	return err
}

func (p *PostgresStorageUnit) Get(ctx context.Context, key string) (string, error) {
	var user User
	err := p.db.QueryRowContext(ctx, "SELECT id, details FROM users WHERE id = $1", key).Scan(&user.ID, &user.Details)
	if err != nil {
		return "", err
	}
	userJson, err := json.Marshal(user)
	if err != nil {
		return "", err
	}
	return string(userJson), nil
}

func (p *PostgresStorageUnit) Delete(ctx context.Context, key string) error {
	_, err := p.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", key)
	return err
}
