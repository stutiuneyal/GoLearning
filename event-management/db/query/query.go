package query

/*

Create Tables

*/

const CreateEventsTable = `

CREATE TABLE IF NOT EXISTS events (
	id SERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	description TEXT NOT NULL,
	location TEXT NOT NULL,
	datetime TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	user_id INTEGER,
	CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
)

`

const InsertEventQuery = `

INSERT INTO events (name, description, location, dateTime, user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING id

`

const GetAllEventsQuery = `

SELECT id,name,description,location,dateTime,user_id FROM events WHERE user_id=$1;

`

const GetEventByIdQuery = `

SELECT id, name, description, location, dateTime, user_id FROM events
WHERE id = $1 and user_id=$2;

`

const UpdateEventQuery = `

UPDATE events SET
                name = $3,
                description = $4,
                location = $5,
                dateTime = $6
WHERE id = $1 and user_id=$2;

`

const DeleteEventQuery = `

DELETE FROM events WHERE id=$1 and user_id=$2;

`

const CreateUsersTable = `

CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	email TEXT NOT NULL,
	password TEXT NOT NULL
)

`

const SignupUserQuery = `

INSERT INTO users (email,password)
VALUES ($1,$2)
RETURNING id

`
const GetUserByEmail = `

SELECT id,email,password
FROM users
WHERE email=$1

`
