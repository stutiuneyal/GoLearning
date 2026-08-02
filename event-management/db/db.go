package database

import (
	"database/sql"
	"fmt"

	"example.com/learning/event-management/db/query"
)

func CreateTables(db *sql.DB) error {

	if _, err := db.Exec(query.CreateUsersTable); err != nil {
		return err
	}

	fmt.Println("users table created")

	if _, err := db.Exec(query.AddUserProfileColumns); err != nil {
		return err
	}

	fmt.Println("user profile columns are ready")

	if _, err := db.Exec(query.CreateUniqueUserEmailIndex); err != nil {
		return err
	}

	fmt.Println("user email unique index is ready")

	if _, err := db.Exec(query.CreateEventsTable); err != nil {
		return err
	}

	fmt.Println("events table created")

	return nil

}
