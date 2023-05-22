package db

import (
	"fmt"
	"os"
	"time"

	"todo/internal/app/constants"
	"todo/internal/app/service/logger"

	"bitbucket.org/liamstask/goose/lib/goose"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var db *gorm.DB
var err error

type DBService struct{}

// Init : Initializes the database migrations
func Init() {

	dbUserName := viper.GetString("DB.engineConfig.user")
	dbPassword := viper.GetString("DB.engineConfig.password")
	dbHost := viper.GetString("DB.engineConfig.host")
	dbName := viper.GetString("DB.engineConfig.dbname")
	dbPort := viper.GetString("DB.engineConfig.port")
	dbURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUserName, dbPassword, dbName)

	maxIdleConnections := viper.GetInt("DB.engineConfig.maxIdleConnection")
	maxOpenConnections := viper.GetInt("DB.engineConfig.maxOpenConnection")
	connectionMaxLifetime := viper.GetInt("DB.engineConfig.connectionMaxLifeTime")

	//dbConnectionString := dbUserName + ":" + dbPassword + "@tcp(" + dbUrl + ")/" + dbName
	db, err = gorm.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("failed to connect.", dbURL, err)
		logger.SugarLogger.Fatalf("Failed to connect to DB", dbURL, err.Error())
		os.Exit(1)
	}

	constants.DBLOGMODE = viper.GetBool("DB.engineConfig.LogMode")
	db.DB().SetMaxIdleConns(maxIdleConnections)
	db.DB().SetMaxOpenConns(maxOpenConnections)
	db.DB().SetConnMaxLifetime(time.Hour * time.Duration(connectionMaxLifetime))
	db.SingularTable(true)

	workingDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Not able to fetch the working directory")
		logger.SugarLogger.Fatalf("Not able to fetch the working directory")
		os.Exit(1)
	}
	workingDir = workingDir + "/internal/app/db/migrations"
	migrateConf := &goose.DBConf{
		MigrationsDir: workingDir,
		Driver: goose.DBDriver{
			Name:    "postgresql",
			OpenStr: dbURL,
			Import:  "github.com/lib/pq",
			Dialect: &goose.PostgresDialect{},
		},
	}
	logger.SugarLogger.Infof("Fetching the most recent DB version")
	latest, err := goose.GetMostRecentDBVersion(migrateConf.MigrationsDir)
	if err != nil {
		logger.SugarLogger.Errorf("Unable to get recent goose db version", err)

	}
	fmt.Println(" Most recent DB version ", latest)
	logger.SugarLogger.Infof("Running the migrations on db", workingDir)
	err = goose.RunMigrationsOnDb(migrateConf, migrateConf.MigrationsDir, latest, db.DB())
	if err != nil {
		logger.SugarLogger.Fatalf("Error while running migrations", err)
		os.Exit(1)
	}
}

// GetDB : Get an instance of DB to connect to the database connection pool
func (d DBService) GetDB() *gorm.DB {
	return db
}
