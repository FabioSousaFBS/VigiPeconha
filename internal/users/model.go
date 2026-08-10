package users

import "time"

type User struct {
	ID             string    `db:"id" json:"id"`
	Name           string    `db:"name" json:"name"`
	Email          string    `db:"email" json:"email"`
	Phone          string    `db:"phone" json:"phone"`
	PasswordHash   string    `db:"password_hash" json:"-"`
	Role           string    `db:"role" json:"role"`
	Status         string    `db:"status" json:"status"`
	OrganizationID *string   `db:"organization_id" json:"organization_id,omitempty"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}
