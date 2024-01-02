package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/mdmdirector/mdmdirector/utils"
	"github.com/pkg/errors"

	// Need to import mysql
	"gorm.io/driver/mysql"
)

var DB *gorm.DB

func Open() error {

	username := utils.DBUsername()
	password := utils.DBPassword()
	dbName := utils.DBName()
	dbHost := utils.DBHost()
	dbPort := utils.DBPort()
	dbSSLMode := utils.DBSSLMode()

	// dbURI := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s", dbHost, dbPort, username, dbName, dbSSLMode, password)
	dbURI := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?tls=%s&charset=utf8&parseTime=True&loc=Local", username, password, dbHost, dbPort, dbName, dbSSLMode)

	var newLogger logger.Interface
	if utils.DebugMode() {
		newLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second, // Slow SQL threshold
				LogLevel:      logger.Info, // Log level
				Colorful:      true,        // Disable color
			},
		)
	} else {
		newLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second,   // Slow SQL threshold
				LogLevel:      logger.Silent, // Log level
				Colorful:      false,         // Disable color
			},
		)
	}

	var err error
	DB, err = gorm.Open(mysql.Open(dbURI), &gorm.Config{Logger: newLogger, DisableForeignKeyConstraintWhenMigrating: true})
	if err != nil {
		return errors.Wrap(err, "Open DB")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return errors.Wrap(err, "creating sqldb object")
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	if utils.DBMaxIdleConnections() != -1 {
		sqlDB.SetMaxIdleConns(utils.DBMaxIdleConnections())
	}

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(utils.DBMaxConnections())

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	return nil
}
