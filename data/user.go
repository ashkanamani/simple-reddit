package data

import (
	"time"
)

type User struct {
	Id        int
	Uuid      string
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
}

type Session struct {
	Id        int
	Uuid      string
	Email     string
	UserId    int
	CreatedAt time.Time
}

// Create a new session for an existing user
func (user *User) CreateSession() (Session, error) {
	var session Session
	statement := "insert into sessions (uuid, email, user_id, created_at) values ($1, $2, $3, $4) returning id, uuid, email, user_id, created_at"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return session, err
	}
	defer stmt.Close()
	// use QueryRow to return a row and scan the returned id into the Session struct
	err = stmt.QueryRow(createUUID(), user.Email, user.Id, time.Now()).Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
	return session, err
}

// Get the session for an existing user
func (user *User) Session() (Session, error) {
	var session Session
	var err error
	session = Session{}
	err = Db.QueryRow("SELECT id, uuid, email, user_id, created_at FROM sessions WHERE user_id = $1", user.Id).
		Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
	return session, err
}

// Check if session is valid in the database
func (session *Session) Check() (bool, error) {
	var valid bool
	err := Db.QueryRow("SELECT id, uuid, email, user_id, created_at FROM sessions WHERE uuid = $1", session.Uuid).
		Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
	if err != nil {
		valid = false
		return valid, err
	}
	if session.Id != 0 {
		valid = true
	}
	return valid, err
}

// Delete session from database
func (session *Session) DeleteByUUID() error {
	statement := "delete from sessions where uuid = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(session.Uuid)
	return err
}

// Get the user from the session
func (session *Session) User() (User, error) {
	user := User{}
	err := Db.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = $1", session.UserId).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedAt)
	return user, err
}

// Delete all sessions from database
func SessionDeleteAll() error {
	statement := "delete from sessions"
	_, err := Db.Exec(statement)
	return err
}

// Create a new user, save user info into the database
func (user *User) Create() error {
	// Postgres does not automatically return the last insert id, because it would be wrong to assume
	// you're always using a sequence.You need to use the RETURNING keyword in your insert to get this
	// information from postgres.
	statement := "insert into users (uuid, name, email, password, created_at) values ($1, $2, $3, $4, $5) returning id, uuid, created_at"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// use QueryRow to return a row and scan the returned id into the User struct
	err = stmt.QueryRow(createUUID(), user.Name, user.Email, Encrypt(user.Password), time.Now()).Scan(&user.Id, &user.Uuid, &user.CreatedAt)
	return err
}

// Delete user from database
func (user *User) Delete() error {
	statement := "delete from users where id = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Id)
	return err
}

// Update user information in the database
func (user *User) Update() error {
	statement := "update users set name = $2, email = $3 where id = $1"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.Id, user.Name, user.Email)
	return err
}

// Delete all users from database
func UserDeleteAll() error {
	statement := "delete from users"
	_, err := Db.Exec(statement)
	return err
}

// Get all users in the database and returns it
func Users() ([]User, error) {

	var users []User
	rows, err := Db.Query("SELECT id, uuid, name, email, password, created_at FROM users")
	if err != nil {
		return users, err
	}
	for rows.Next() {
		user := User{}
		if err = rows.Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt); err != nil {
			return users, err
		}
		users = append(users, user)
	}
	rows.Close()
	return users, err
}

// Get a single user given the email
func UserByEmail(email string) (User, error) {
	user := User{}
	err := Db.QueryRow("SELECT id, uuid, name, email, password, created_at FROM users WHERE email = $1", email).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	return user, err
}

// Get a single user given the UUID
func UserByUUID(uuid string) (User, error) {
	user := User{}
	err := Db.QueryRow("SELECT id, uuid, name, email, password, created_at FROM users WHERE uuid = $1", uuid).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	return user, err
}
