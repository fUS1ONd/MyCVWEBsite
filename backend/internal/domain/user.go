// Package domain contains core business entities and types
package domain

import "time"

// Role represents user role in the system
type Role string

const (
	// RoleAdmin represents administrator role with full access
	RoleAdmin Role = "admin"
	// RoleUser represents regular user role
	RoleUser Role = "user"
)

// User represents a user in the system
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}
