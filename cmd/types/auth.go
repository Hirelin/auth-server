package types

import "github.com/jackc/pgx/v5/pgtype"

type contextKey string

const (
	SessionKey contextKey = "session"
	UserKey    contextKey = "user"
)

type SessionData struct {
	Token pgtype.UUID `json:"token"`
	Valid bool        `json:"valid"`
}

type UserData struct {
	Valid    bool        `json:"valid"`
	ID       pgtype.UUID `json:"id"`
	Name     string      `json:"name"`
	Email    string      `json:"email"`
	UserName string      `json:"username"`
}
