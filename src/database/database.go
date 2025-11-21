package database

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Database struct {
	Host       string
	Port       int
	User       string
	Password   string
	Name       string
	connection *pgxpool.Pool
}

func NewDatabase(host string, port int, user, password, name string) *Database {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", user, password, host, port, name)
	fmt.Println("Database DSN:", dsn) // For demonstration purposes

	connection, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to database: %v\n", err))
	}

	return &Database{
		Host:       host,
		Port:       port,
		User:       user,
		Password:   password,
		Name:       name,
		connection: connection,
	}
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

func (db *Database) QueryWithScan(query string, dest interface{}) error {
	if reflect.ValueOf(dest).Kind() != reflect.Ptr {
		return errors.New("destination should be pointer or else changes won't reflect on it")
	}

	rows, err := db.connection.Query(context.Background(), query)
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
