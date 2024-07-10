package postgresql

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Database struct {
	db *sqlx.DB
}

func NewDatabase(user, password, host, port, dbName string) (*sqlx.DB, error) {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbName)

	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	// Ping базы данных для проверки подключения
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (d *Database) Close() error {
	err := d.db.Close()
	return err
}
