package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/crypto/bcrypt"
)

const (
	connString = "postgres://kubepilot:kubepilot_dev_password@localhost:5432/users?sslmode=disable"
)

type User struct {
	Id             int
	Name           string
	Email          string
	HashedPassword string
	CreatedAt      time.Time
	Profile        Profile
}

type Profile struct {
	UserId    int
	Avatar    string
	CreatedAt time.Time
}

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

	db, err := sql.Open("pgx", connString)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Successfully connected to the database")

	// create the tables
	createTables(db)

	// write queries
	var id int64
	if id, err = createUser(db, "Sparry", "sparry@gmail.com", "sparry@123", "hhtps://avatar.com/sparry"); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("User created: ", id)
	}

	if user, err := getUserById(db, id); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Printf("User: %+v\n", user)
	}

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

// Tx: Begin, Rollback, Commit
// Same Tx -> spans across multiple queries -> same context
func createUser(db *sql.DB, name, email, password, avatar string) (int64, error) {

	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return 0, err
	}

	userInsertQuery := `
		INSERT INTO users (name, email, hashed_password)
		VALUES ($1,$2,$3)
		RETURNING id;
	`

	user := User{
		Name:           name,
		Email:          email,
		HashedPassword: string(hashedPassword),
	}

	if err := db.QueryRow(userInsertQuery, user.Name, user.Email, user.HashedPassword).Scan(&user.Id); err != nil {
		return 0, err
	}

	// user is created successfully -> proceed and create user profile

	userProfileInsertQuery := `
		INSERT INTO profiles (user_id,avatar)
		VALUES ($1,$2)
	`

	if _, err := db.Exec(userProfileInsertQuery, user.Id, avatar); err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return int64(user.Id), nil

}

func getUserById(db *sql.DB, id int64) (User, error) {

	getUserByIdQuery := `
		SELECT u.id, u.name, u.email, u.hashed_password, p.avatar
		FROM users u INNER JOIN profiles p
		ON u.id = p.user_id
		WHERE u.id=$1;
	`

	var user User

	if err := db.QueryRow(getUserByIdQuery, id).Scan(&user.Id, &user.Name, &user.Email, &user.HashedPassword, &user.Profile.Avatar); err != nil {
		return User{}, err
	}

	return user, nil

}
