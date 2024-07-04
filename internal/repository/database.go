package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"os"
	"sync"
	"wb-weather/pkg/logger"
)

type Database struct {
	db    *sqlx.DB
	mutex sync.Mutex
}

var (
	instance *Database
	once     sync.Once
)

func GetInstance() (*Database, error) {
	var err error
	once.Do(func() {
		logger.Info("Создание нового экземпляра Database")
		instance, err = newDatabase()
		if err != nil {
			logger.Fatal("Ошибка создания экземпляра Database", zap.Error(err))
		}
	})
	return instance, err
}

func newDatabase() (*Database, error) {
	logger.Info("Загрузка файла .env для конфигурации базы данных")
	err := godotenv.Load()
	if err != nil {
		logger.Fatal("Ошибка загрузки файла .env", zap.Error(err))
	}

	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbName := os.Getenv("POSTGRES_DB")

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

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

	return &Database{db: db}, nil
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
