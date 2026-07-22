package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/crypto/bcrypt"
)

const (
	connString = "postgres://kubepilot:kubepilot_dev_password@localhost:5432/users?sslmode=disable"
)

type User struct {
	Name           string
	Email          string
	HashedPassword string
}

var (
	insertUserStmt     *sql.Stmt
	getUserByIdStmt    *sql.Stmt
	getUserByEmailStmt *sql.Stmt
	getAllUsersStmt    *sql.Stmt
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

	// Skipping create table -> that is already done

	/*
		Prepare the statments once
	*/
	insertUserStmt, err = db.Prepare(
		`
			INSERT INTO users (name, email, hashed_password)
			VALUES ($1,$2,$3)
			RETURNING id;
		`,
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer insertUserStmt.Close()

	getUserByIdStmt, err = db.Prepare(
		`
			SELECT name, email, hashed_password
			FROM users
			WHERE id=$1;
		`,
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer getUserByIdStmt.Close()

	getUserByEmailStmt, err = db.Prepare(
		`
			SELECT name, email, hashed_password
			FROM users
			WHERE email=$1;
		`,
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer getUserByEmailStmt.Close()

	getAllUsersStmt, err = db.Prepare(
		`
			SELECT name,email,hashed_password
			FROM users;
		`,
	)
	if err != nil {
		fmt.Println(err)
	}
	defer getAllUsersStmt.Close()

	// call the functions
	if id, err := createdUserWithPrepare("Mark", "mark@gmail.com", "1234"); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("User created: ", id)
	}

	if user, err := getUserByIdWithPrepare(1); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Printf("User: %+v\n", user)
	}

	if users, err := getAllUsersWithPrepare(); err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Printf("users: %#v\n", users)
	}

}

func createdUserWithPrepare(name, email, password string) (int64, error) {

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

	var id int64

	if err := insertUserStmt.QueryRow(user.Name, user.Email, user.HashedPassword).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil

}

func getUserByIdWithPrepare(id int64) (User, error) {

	var user User

	if err := getUserByIdStmt.QueryRow(id).Scan(&user.Name, &user.Email, &user.HashedPassword); err != nil {
		return User{}, err
	}

	return user, nil

}

func getAllUsersWithPrepare() ([]User, error) {

	users := []User{}

	rows, err := getAllUsersStmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		var user User

		if err := rows.Scan(&user.Name, &user.Email, &user.HashedPassword); err != nil {
			return nil, err
		}

		users = append(users, user)

	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil

}
