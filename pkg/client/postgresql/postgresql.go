package postgresql

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"wb-weather/pkg/logger"
)

type Database struct {
	db *sqlx.DB
}

func NewDatabase(user, password, host, port, dbName string) (*sqlx.DB, error) {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbName)

	logger.Debug("Строка подключения к базе данных сформирована", zap.String("connectionString", connectionString))

	db, err := sqlx.Connect("postgres", connectionString)
	if err != nil {
		logger.Fatal("Ошибка подключения к базе данных", zap.Error(err))
	}

	// Ping базы данных для проверки подключения
	logger.Info("Проверка подключения к базе данных")
	err = db.Ping()
	if err != nil {
		logger.Fatal("Ошибка проверки подключения к базе данных", zap.Error(err))
	}
	logger.Info("Подключение к базе данных PostgreSQL успешно")

	return db, nil
}

func (d *Database) Close() error {
	logger.Info("Закрытие подключения к базе данных")
	err := d.db.Close()
	if err != nil {
		logger.Error("Ошибка при закрытии подключения к базе данных", zap.Error(err))
	} else {
		logger.Info("Подключение к базе данных успешно закрыто")
	}
	return err
}
