package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DBConfig ... Deklarasi atribut konfigurasi database
type DBConfig struct {
	Host     string
	Port     string
	User     string
	DBName   string
	Password string
}

// DBModel ... Referensi database yang digunakan
type DBModel struct {
	// DB ... Referensi database yang digunakan
	DB gorm.DB
}

// DBModelInterface ... Deklarasi interfaces
type DBModelInterface interface {
	BuildDBConfig() *DBConfig
	DbURL(dbConfig *DBConfig) string
	MysqlConn() (*gorm.DB, error)
	PostgreConn() (*gorm.DB, error)
}

// BuildDBConfig ... Inisialisasi konfigurasi pada database
func (dbsql *DBModel) BuildDBConfig() *DBConfig {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return &DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}
}

// DbMySQL ... Mengambil url yang digunakan untuk driver MySQL
func (dbsql *DBModel) DbMySQL(dbConfig *DBConfig) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName)
}

// DbPostgre ... Mengambil url yang digunakan untuk driver MySQL
func (dbsql *DBModel) DbPostgre(dbConfig *DBConfig) string {
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.DBName,
		dbConfig.Host,
		dbConfig.Port,
	)
}

// MysqlConn ... Mysql Connection
func (dbsql *DBModel) MysqlConn() *gorm.DB {
	result, err := gorm.Open(mysql.Open(dbsql.DbMySQL(dbsql.BuildDBConfig())), &gorm.Config{})
	if err != nil {
		return nil
	}
	return result
}

// PostgreConn ... Postgre Connection
func (dbsql *DBModel) PostgreConn() *gorm.DB {
	result, err := gorm.Open(mysql.Open(dbsql.DbPostgre(dbsql.BuildDBConfig())), &gorm.Config{})
	if err != nil {
		return nil
	}
	return result
}
