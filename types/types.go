package types

import "time"

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

type Action struct {
	ID        int       `json:"id"`
	Type      string    `json:"type"` // use type
	UserID    int       `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}
