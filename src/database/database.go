package database

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/luk3skyw4lker/order-pack-calculator/src/config"
)

type Database struct {
	connection *pgxpool.Pool
}

func NewDatabase(cfg config.DatabaseConfig) (*Database, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	fmt.Println("Database DSN:", dsn) // For demonstration purposes

	connection, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to database: %v\n", err))
	}

	return &Database{
		connection: connection,
	}, err
}

func (c *Database) scanValue(rows pgx.Rows, dest interface{}) error {
	valueKind := reflect.ValueOf(dest).Elem().Kind()

	switch valueKind {
	case reflect.Slice:
		return pgxscan.ScanAll(dest, rows)
	default:
		return pgxscan.ScanOne(dest, rows)
	}
}

func (db *Database) QueryWithScan(query string, dest interface{}, args ...any) error {
	if reflect.ValueOf(dest).Kind() != reflect.Ptr {
		return errors.New("destination should be pointer or else changes won't reflect on it")
	}

	rows, err := db.connection.Query(context.Background(), query, args...)
	if err != nil {
		return err
	}

	return db.scanValue(rows, dest)
}

func (db *Database) Query(query string) error {
	_, err := db.connection.Query(context.Background(), query)
	return err
}

func (db *Database) Close() {
	db.connection.Close()
}
