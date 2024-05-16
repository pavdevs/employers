package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type Database struct {
	db     *sql.DB
	config Config
	logger *logrus.Logger
}

func NewDatabase(config Config, logger *logrus.Logger) *Database {
	return &Database{
		config: config,
		logger: logger,
	}
}

func (d *Database) Connect() error {
	connectionString := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=%s password=%s", d.config.Host, d.config.User, d.config.DBName, d.config.SSLMode, d.config.Password)
	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		d.logger.Error(err)
		return err
	}

	d.logger.Info("Database prepared")

	if err := db.Ping(); err != nil {
		d.logger.Error(err)
		return err
	}

	d.db = db

	d.logger.Info("Database connected")

	return nil
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}

func (d *Database) Disconnect() error {
	if err := d.db.Close(); err != nil {
		d.logger.Error(err)
		return err
	}

	d.logger.Info("Database disconnected")
	return nil
}
