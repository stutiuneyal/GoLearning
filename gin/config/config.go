package config

import (
	"database/sql"

	"example.com/learning/gin/constants"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnectToDatabase() (*sql.DB, error) {

	db, err := sql.Open("pgx", constants.DbConnString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	return db, nil

}
