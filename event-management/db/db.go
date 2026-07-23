package database

import (
	"database/sql"
	"fmt"

	"example.com/learning/gin/db/query"
)

func CreateTables(db *sql.DB) error {

	if _, err := db.Exec(query.CreateUsersTable); err != nil {
		return err
	}

	fmt.Println("users table created")

	if _, err := db.Exec(query.CreateEventsTable); err != nil {
		return err
	}

	fmt.Println("events table created")

	return nil

}
