package db

// schema.go provides data models in DB
import (
	"time"
)

type User struct {
    ID        uint64    `db:"id"`
    Name      string    `db:"name"`
    Password  []byte    `db:"password"`
	Updated_At time.Time `db:"updated_at"`
	Created_At time.Time `db:"created_at"`
}

// Task corresponds to a row in `tasks` table
type Task struct {
	ID        uint64    `db:"id"`
	Title     string    `db:"title"`
	CreatedAt time.Time `db:"created_at"`
	IsDone    bool      `db:"is_done"`
	Detail    string    `db:"detail"`
}

type Ownership struct {
	User_id   uint64   `db:"user_id"`
	Task_id   uint64   `db:"task_id"`
}