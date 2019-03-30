package config

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

var (
	DBUSER     = "postgres"
	DBPASSWORD = "2017"
	DBNAME     = "tp"

	CONNBDStr         = " user=" + DBUSER + " dbname=" + DBNAME + " password=" + DBPASSWORD + " sslmode=disable"
	DB        *sql.DB = connectDB(CONNBDStr)
	PORT      string  = "8080"
)

func connectDB(CONNBDStr string) *sql.DB {
	DB, err := sql.Open("postgres", CONNBDStr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Database connected!")

	return DB
}
