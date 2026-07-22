package main

import (
	"database/sql"
	"fmt"

	"example.com/learning/database/4-repository/repository/impl"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	connString = "postgres://kubepilot:kubepilot_dev_password@localhost:5432/users?sslmode=disable"
)

var (
	createUserTableQuery = `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL,
			hashed_password BYTEA NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`

	createUserProfileQuery = `
		CREATE TABLE IF NOT EXISTS profiles(
			user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
			avatar TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		)
	`
)

func main() {

	db, err := connectDatabase()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	fmt.Println("Successfully connected to the database")

	// create the tables
	createTables(db)

	// initialize my repository
	userRepository := impl.NewUserRepositoryImpl(db)

	// write queries
	var id int64
	if id, err = userRepository.CreateUser("Harry", "harry@gmail.com", "harry@123", "hhtps://avatar.com/harry"); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("User created: ", id)
	}

	if user, err := userRepository.GetUserById(id); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Printf("User: %+v\n", user)
	}

}

func connectDatabase() (*sql.DB, error) {

	db, err := sql.Open("pgx", connString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil

}

func createTables(db *sql.DB) {

	if _, err := db.Exec(createUserTableQuery); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("users table created successfully")

	if _, err := db.Exec(createUserProfileQuery); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("User Profile table created successfully")

}
