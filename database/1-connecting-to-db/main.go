package main

import (
	"database/sql"
	"errors"
	"fmt"

	// blank import beause we only want the init() function to run, which runs as soon as the package is imported
	// this init() function then registers the postgres(pgx) driver inside "database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/crypto/bcrypt"
)

const (
	connString = "postgres://kubepilot:kubepilot_dev_password@localhost:5432/users?sslmode=disable"
)

var (
	createUserTableSchema = `
	CREATE TABLE IF NOT EXISTS users(
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL,
		hashed_password BYTEA NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`
)

type User struct {
	Name           string
	Email          string
	HashedPassword string
}

func main() {

	// Create a db connection
	db, err := sql.Open("pgx", connString)
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		fmt.Println("Closing the database")
		if err := db.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// verify the connection
	if err := db.Ping(); err != nil {
		fmt.Println("Error pinging database ", err)
	}

	fmt.Println("Successfully connected to database")

	createTables(db)

	// Insert Users
	if found, err := GetUserByEmail(db, "john@example.com"); err != nil {
		fmt.Println("User not found: ", err)
	} else if !found {
		if id, err := createUser(db, "John", "jogn@example.com", "1234"); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("User Created: ", id)
		}
	} else if found {
		fmt.Println("User already exists")
	}

	// get user by id
	if user, err := GetUserById(db, 1); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%+v\n", user)
	}

	// get all users
	if users, err := ListUsers(db); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%#v\n", users)
	}

}

func createTables(db *sql.DB) {

	if _, err := db.Exec(createUserTableSchema); err != nil {
		fmt.Printf("Error creating table: %v\n", err)
	}

	fmt.Println("Successfully created users table")

}

func createUser(db *sql.DB, name, email, password string) (int64, error) {

	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return 0, err
	}

	user := User{
		Name:           name,
		Email:          email,
		HashedPassword: string(hashedPassword),
	}

	userInsertQuery := `
		INSERT INTO users (name, email, hashed_password)
		VALUES ($1,$2,$3)
		RETURNING id;
	`

	var id int64

	if err := db.QueryRow(userInsertQuery, user.Name, user.Email, user.HashedPassword).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil

}

func GetUserById(db *sql.DB, id int64) (User, error) {

	getUserByIdQuery := `
		SELECT name, email, hashed_password
		FROM users
		WHERE id=$1
	`

	var user User

	if err := db.QueryRow(getUserByIdQuery, id).Scan(&user.Name, &user.Email, &user.HashedPassword); err != nil {
		return User{}, err
	}

	return user, nil

}

func GetUserByEmail(db *sql.DB, email string) (bool, error) {

	getUserByIdQuery := `
		SELECT name, email, hashed_password
		FROM users
		WHERE email=$1
	`

	var user User

	if err := db.QueryRow(getUserByIdQuery, email).Scan(&user.Name, &user.Email, &user.HashedPassword); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil // User does not exist
		}
		return false, err // Some other database error
	}

	return true, nil // user already exists

}

func ListUsers(db *sql.DB) ([]User, error) {

	getAllUsersQuery := `
		SELECT name, email, hashed_password
		FROM users
	`

	rows, err := db.Query(getAllUsersQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User

	for rows.Next() {

		var user User

		if err := rows.Scan(&user.Name, &user.Email, &user.HashedPassword); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	// Important: check for iteration error
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil

}
