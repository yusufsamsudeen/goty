package goty

import (
	"errors"
	"fmt"
	"github.com/olebedev/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io/ioutil"
	"sync"
	"time"
)

func loadConfig() (*config.Config,error) {

	file,err := ioutil.ReadFile("config.yml")
	if err!=nil {
		return nil, errors.New("no config file found")
	}

	yamlString := string(file)
	cfg, err := config.ParseYaml(yamlString)

	return cfg,err
}

func getConfiguration(path string) *config.Config {
	cfg, err := loadConfig()
	if err!=nil {
		fmt.Println(err)
		return nil
	}

	cfg, _ = cfg.Get(path)
	return cfg
}

type connectionStruct struct {
	db *gorm.DB
	config *connectionConfig
}

type connectionConfig struct {
	driver string
	host string
	port string
	username string
	password string
	database string
	migrate bool
}

var connection *connectionStruct
var once sync.Once

func getConnectionString() (*connectionConfig, error) {
	databaseConfig := getConfiguration("database")
	if databaseConfig==nil {
		return nil, errors.New("no database config found")
	}

	driver, _ := databaseConfig.String("driver")
	username, _ := databaseConfig.String("username")
	password, _ := databaseConfig.String("password")
	database, _ := databaseConfig.String("schema")
	migrate, _ := databaseConfig.Bool("migrate")
	host, _ := databaseConfig.String("host")
	port, _ := databaseConfig.String("port")

	return &connectionConfig{
		driver:   driver,
		username: username,
		password: password,
		database: database,
		migrate:  migrate,
		host: host,
		port: port,
	}, nil
}

func acquireConnection() (*gorm.DB, *connectionConfig) {

	connectionString, err := getConnectionString()
	if err != nil {
		fmt.Println(err)
		return nil,nil
	}

	driver := connectionString.driver
	switch driver {
	case "mysql":
		return mysqlConnection(connectionString), connectionString
	case "postgres":
		return postgresConnection(connectionString), connectionString

	default:
		return nil,nil
	}



}

func mysqlConnection(connectionString *connectionConfig) *gorm.DB{

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4", connectionString.username, connectionString.password,
		connectionString.host, connectionString.port, connectionString.database)
	db, err  := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err !=nil {
		panic("could not establish connection")
	}


	sqlDB, err := db.DB()
	if err!=nil {
		panic("error establishing connection")
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db
}

func postgresConnection(config *connectionConfig) *gorm.DB{


	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", config.host, config.username, config.password, config.database, config.port)
	db, err  := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err !=nil {
		panic("could not establish connection")
	}


	sqlDB, err := db.DB()
	if err!=nil {
		panic("error establishing connection")
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db
}

func getConnection() (*gorm.DB, *connectionConfig) {
	once.Do(func() {
		db, connectionConfig := acquireConnection()
		connection = &connectionStruct{db : db, config: connectionConfig}
	})
	return connection.db, connection.config
}

func Open() *gorm.DB {
	connection, _ := getConnection()
	return connection
}

func Close() bool {
	sqlDB, err := connection.db.DB()
	if err==nil {
		panic("Could not communicate with DB")
	}
	err = sqlDB.Close()
	if err != nil {
		panic("Could not close DB")
	}

	return true
}

