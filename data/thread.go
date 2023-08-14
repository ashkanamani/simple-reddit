package data

import (
	"time"
)

type Thread struct {
	Id        int
	Uuid      string
	Topic     string
	UserId    int
	CreatedAt time.Time
}

type Post struct {
	Id        int
	Uuid      string
	Body      string
	UserId    int
	ThreadId  int
	CreatedAt time.Time
}

// format the CreatedAt date to display nicely on the screen
func (thread *Thread) CreatedAtDate() string {
	return thread.CreatedAt.Format("Jan 2, 2006 at 3:04pm")
}

func (post *Post) CreatedAtDate() string {
	return post.CreatedAt.Format("Jan 2, 2006 at 3:04pm")
}

// get the number of posts in a thread
func (thread *Thread) NumReplies() int {
	var count int
	rows, err := Db.Query("SELECT count(*) FROM posts where thread_id = $1", thread.Id)
	if err != nil {
		return count
	}
	for rows.Next() {
		if err = rows.Scan(&count); err != nil {
			return count
		}
	}
	rows.Close()
	return count
}

// get posts to a thread
func (thread *Thread) Posts() ([]Post, error) {
	var posts []Post
	rows, err := Db.Query("SELECT id, uuid, body, user_id, thread_id, created_at FROM posts where thread_id = $1", thread.Id)
	if err != nil {
		return posts, err
	}
	for rows.Next() {
		post := Post{}
		if err = rows.Scan(&post.Id, &post.Uuid, &post.Body, &post.UserId, &post.ThreadId, &post.CreatedAt); err != nil {
			return posts, err
		}
		posts = append(posts, post)
	}
	rows.Close()
	return posts, nil
}

// Create a new thread
func (user *User) CreateThread(topic string) (Thread,error) {
	var conv Thread
	statement := "insert into threads (uuid, topic, user_id, created_at) values ($1, $2, $3, $4) returning id, uuid, topic, user_id, created_at"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return conv, err
	}
	defer stmt.Close()
	// use QueryRow to return a row and scan the returned id into the Session struct
	err = stmt.QueryRow(createUUID(), topic, user.Id, time.Now()).Scan(&conv.Id, &conv.Uuid, &conv.Topic, &conv.UserId, &conv.CreatedAt)
	return conv, err
}

// Create a new post to a thread
func (user *User) CreatePost(conv Thread, body string) (Post, error) {
	var post Post
	var err error
	statement := "insert into posts (uuid, body, user_id, thread_id, created_at) values ($1, $2, $3, $4, $5) returning id, uuid, body, user_id, thread_id, created_at"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return post, err
	}
	defer stmt.Close()
	// use QueryRow to return a row and scan the returned id into the Session struct
	err = stmt.QueryRow(createUUID(), body, user.Id, conv.Id, time.Now()).Scan(&post.Id, &post.Uuid, &post.Body, &post.UserId, &post.ThreadId, &post.CreatedAt)
	return post, err
}

// Get all threads in the database and returns it
func Threads() ([]Thread, error) {
	var threads []Thread
	var err error
	rows, err := Db.Query("SELECT id, uuid, topic, user_id, created_at FROM threads ORDER BY created_at DESC")
	if err != nil {
		return threads, err
	}
	for rows.Next() {
		conv := Thread{}
		if err = rows.Scan(&conv.Id, &conv.Uuid, &conv.Topic, &conv.UserId, &conv.CreatedAt); err != nil {
			return threads, err
		}
		threads = append(threads, conv)
	}
	rows.Close()
	return threads, err
}

// Get a thread by the UUID
func ThreadByUUID(uuid string) (Thread, error) {
	var conv Thread 
	var err error
	conv = Thread{}
	err = Db.QueryRow("SELECT id, uuid, topic, user_id, created_at FROM threads WHERE uuid = $1", uuid).
		Scan(&conv.Id, &conv.Uuid, &conv.Topic, &conv.UserId, &conv.CreatedAt)
	return conv, err
}

// Get the user who started this thread
func (thread *Thread) User() User {
	user := User{}
	Db.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = $1", thread.UserId).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedAt)
	return user
}

// Get the user who wrote the post
func (post *Post) User() User {
	user := User{}
	Db.QueryRow("SELECT id, uuid, name, email, created_at FROM users WHERE id = $1", post.UserId).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedAt)
	return user
}
