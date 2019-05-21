package config

import (
	"github.com/jackc/pgx"
)

var (
	// DBUSER     = "postgres"
	// DBPASSWORD = "2017"
	// DBNAME     = "highload"
	DBUSER     = "docker"
	DBPASSWORD = "docker"
	DBNAME     = "docker"

	CONNBDStr        = " user=" + DBUSER + " dbname=" + DBNAME + " password=" + DBPASSWORD + " sslmode=disable"
	DB               = connectBDpgx(config)
	PORT      string = "5000"
	// PORT      string  = "5000"
)

// func connectDB(CONNBDStr string) *sql.DB {
// 	DB, err := sql.Open("postgres", CONNBDStr)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	log.Println("Database connected!")

// 	return DB
// }

var config = pgx.ConnPoolConfig{
	ConnConfig: pgx.ConnConfig{
		Host:     "localhost",
		Port:     5432,
		User:     DBUSER,
		Password: DBPASSWORD,
		Database: DBNAME,
	},
	MaxConnections: 30,
}

func connectBDpgx(config pgx.ConnPoolConfig) *pgx.ConnPool {
	pool, err := pgx.NewConnPool(config)
	if err != nil {
		Logger.Fatal(err.Error())
	}

	return pool
}
