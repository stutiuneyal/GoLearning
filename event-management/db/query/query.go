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
	password TEXT NOT NULL,
	name VARCHAR(100) NOT NULL DEFAULT '',
	bio VARCHAR(500) NOT NULL DEFAULT '',
	profile_picture_key TEXT
)

`

const AddUserProfileColumns = `

ALTER TABLE users
	ADD COLUMN IF NOT EXISTS name VARCHAR(100) NOT NULL DEFAULT '',
	ADD COLUMN IF NOT EXISTS bio VARCHAR(500) NOT NULL DEFAULT '',
	ADD COLUMN IF NOT EXISTS profile_picture_key TEXT

`

const CreateUniqueUserEmailIndex = `

	CREATE UNIQUE INDEX IF NOT EXISTS users_email_unique_idx
	ON users(email)

`

const SignupUserQuery = `

INSERT INTO users (email,password,name,bio)
VALUES ($1,$2,$3,$4)
RETURNING id

`
const GetUserByEmail = `

SELECT id,email,password,name,bio,COALESCE(profile_picture_key,'')
FROM users
WHERE email=$1

`

const GetProfileByUserIDQuery = `

	SELECT id,email,name,bio,COALESCE(profile_picture_key,'')
	FROM users
	WHERE id=$1

`

const UpdateProfileQuery = `

	UPDATE users
	SET
		name = CASE
			WHEN $2 THEN $3
			ELSE name
		END,
		bio = CASE
			WHEN $4 THEN $5
			ELSE bio
		END
	WHERE id = $1
	RETURNING id,email,name,bio,COALESCE(profile_picture_key,'')

`

/*
FOR UPDATE
Locks the selected row for updating.
Other transactions cannot modify or acquire another FOR UPDATE lock on the same row until the current transaction commits or rolls back.
This helps prevent race conditions when multiple transactions try to update the same user.
*/
const GetProfilePictureForUpdateQuery = `

SELECT
	COALESCE(profile_picture_key, '')
FROM users
WHERE id = $1
FOR UPDATE

`

const UpdateProfilePictureQuery = `

UPDATE users
SET profile_picture_key = $2
WHERE id = $1

`

const RemoveProfilePictureQuery = `

UPDATE users
SET profile_picture_key = NULL
WHERE id = $1

`
