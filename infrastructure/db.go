package infrastructure

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/sayopaul/sendchamp-go-test/config"
	"github.com/sayopaul/sendchamp-go-test/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	*gorm.DB
	dsn string
}

type DatabaseCred struct {
	DBUsername string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	SSLMode    string
}

// support database string for db initialization
func getDatabaseCred(databaseUrl string) DatabaseCred {
	str := strings.Split(databaseUrl, "@")
	part1 := str[0]
	part2 := str[1]

	splitPart2 := strings.Split(part2, ":")

	host := splitPart2[0]

	part3 := splitPart2[1]
	splitPart3 := strings.Split(part3, "/")

	port := splitPart3[0]
	dbname := splitPart3[1]

	sslMode := "disable"
	if strings.Contains(dbname, "?") {
		spiltPart5 := strings.Split(dbname, "?")
		dbname = spiltPart5[0]
		sslModePart := spiltPart5[1]
		sslMode = strings.Split(sslModePart, "=")[1]
	}

	splitPart1 := strings.Split(part1, ":")

	password := splitPart1[2]

	part4 := splitPart1[1]
	splitPart4 := strings.Split(part4, "//")

	username := splitPart4[1]
	return DatabaseCred{
		DBUsername: username,
		DBPassword: password,
		DBHost:     host,
		DBPort:     port,
		DBName:     dbname,
		SSLMode:    sslMode,
	}
}

// NewDatabase creates a new database instance
func NewDatabase(configEnv config.Config) Database {

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	dbCred := DatabaseCred{
		DBUsername: configEnv.DBUsername,
		DBPassword: configEnv.DBPassword,
		DBHost:     configEnv.DBHost,
		DBPort:     configEnv.DBPort,
		DBName:     configEnv.DBName,
		SSLMode:    configEnv.DBSslMode,
	}

	if configEnv.DatabaseUrl != "" {
		dbCred = getDatabaseCred(configEnv.DatabaseUrl)
	}

	// "gorm:gorm@tcp(localhost:9910)/gorm?charset=utf8&parseTime=True&loc=Local"
	dbUrl := fmt.Sprintf("%v:%v/@tcp(%v:%v)/%v",
		dbCred.DBUsername, dbCred.DBPassword, dbCred.DBHost, dbCred.DBPort, dbCred.DBName)

	db, err := gorm.Open(mysql.Open(dbUrl), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Println("Url: ", dbUrl)
		log.Panic(err)
	}

	log.Println("Database connection established")

	//automigrate tables
	err = db.AutoMigrate(
		&models.User{},
	)
	if err != nil {
		log.Panic(err)
	}

	return Database{
		DB:  db,
		dsn: dbUrl,
	}
}
