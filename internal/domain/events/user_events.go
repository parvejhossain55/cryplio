package events

import "time"

// UserRegisteredEvent is raised when a new user registers.
type UserRegisteredEvent struct {
	UserID    string    `json:"user_id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

func (e UserRegisteredEvent) Name() string { return "user.registered" }

// UserLoggedInEvent is raised when a user logs in.
type UserLoggedInEvent struct {
	UserID     string    `json:"user_id"`
	Email      string    `json:"email"`
	LoggedInAt time.Time `json:"logged_in_at"`
}

func (e UserLoggedInEvent) Name() string { return "user.logged_in" }

// UserLoggedOutEvent is raised when a user logs out.
type UserLoggedOutEvent struct {
	UserID      string    `json:"user_id"`
	LoggedOutAt time.Time `json:"logged_out_at"`
}

func (e UserLoggedOutEvent) Name() string { return "user.logged_out" }
