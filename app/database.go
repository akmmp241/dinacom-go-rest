package app

import (
	"akmmp241/dinamcom-2024/dinacom-go-rest/config"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func NewDB(cnf *config.Config) *sql.DB {
	DbName := cnf.Env.GetString("DB_NAME")
	DbUser := cnf.Env.GetString("DB_USER")
	DbPass := cnf.Env.GetString("DB_PASS")
	DbHost := cnf.Env.GetString("DB_HOST")
	DbPort := cnf.Env.GetString("DB_PORT")

	conn, err := sql.Open("mysql", DbUser+":"+DbPass+"@tcp("+DbHost+":"+DbPort+")/"+DbName)
	if err != nil {
		log.Fatal("error while connect to database", err)
	}

	err = conn.Ping()
	if err != nil {
		log.Fatal("error while ping to database", err)
	}
	log.Println("Connected to database")

	//conn.SetMaxIdleConns(10)
	//conn.SetConnMaxLifetime(5)
	//conn.SetMaxOpenConns(10)
	//conn.SetConnMaxIdleTime(10)

	return conn
}
