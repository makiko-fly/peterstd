package datasource

import (
	"fmt"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type MySQLConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DBName   string `yaml:"db_name"`
	MaxIdle  int    `yaml:"max_idle"`
	MaxConn  int    `yaml:"max_conn"`
	LogMode  bool   `yaml:"log_mode"`
}

func (c MySQLConfig) New() *gorm.DB {
	return NewMySQLClient(
		c.Host,
		c.Port,
		c.Username,
		c.Password,
		c.DBName,
		c.MaxIdle,
		c.MaxConn,
		c.LogMode,
	)
}

func NewMySQLClient(host string, port int, username, password, dbname string, maxIdle, maxConn int, log bool) *gorm.DB {
	if maxIdle == 0 {
		maxIdle = 50
	}
	if maxConn == 0 {
		maxConn = 100
	}
	dbURL := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=%s",
		username,
		password,
		host,
		port,
		dbname,
		url.QueryEscape("Asia/Shanghai"),
	)

	fmt.Printf("Try to connect to MYSQL %s:%d\n", host, port)
	db, err := gorm.Open("mysql", dbURL)
	if err != nil {
		panic(fmt.Sprintf("failed to connect MYSQL %s, %s", dbURL, err.Error()))
	}
	fmt.Println("Connected to MYSQL", host, ":", port, ", logMode: ", log)
	db.LogMode(log)
	db.SingularTable(true)
	db.DB().SetMaxIdleConns(maxIdle)
	db.DB().SetMaxOpenConns(maxConn)
	db.DB().SetConnMaxLifetime(time.Hour)
	db.AutoMigrate()
	return db
}
