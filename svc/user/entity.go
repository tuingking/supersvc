package user

import (
	"database/sql"
	"time"
)

const (
	UserStatusInActive = iota
	UserStatusActive
	UserStatusBanned
	UserStatusSuspend
)

type User struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Phone     string    `json:"phone" db:"phone"`
	Email     string    `json:"email" db:"email"`
	Status    int       `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type CreateUser struct {
	Name   string `json:"name" db:"name"`
	Phone  string `json:"phone" db:"phone"`
	Email  string `json:"email" db:"email"`
	Status int    `json:"status" db:"status"`
}

type GetUserParam struct {
	Email sql.NullString `param:"email" db:"email"`

	Page   int64    `param:"page"`
	Limit  int64    `param:"limit"`
	SortBy []string `param:"sortBy"`
}
