// Package domain contains core business entities and types
package domain

import "time"

// Contacts represents contact information
type Contacts struct {
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
	GitHub   string `json:"github,omitempty" validate:"omitempty,url"`
	LinkedIn string `json:"linkedin,omitempty" validate:"omitempty,url"`
	VK       string `json:"vk,omitempty" validate:"omitempty,url"`
}

// Profile represents CV/profile information for the website owner
type Profile struct {
	ID          int       `json:"id"`
	Name        string    `json:"name" validate:"required,min=1,max=255"`
	Description string    `json:"description" validate:"required,min=1"`
	PhotoURL    string    `json:"photo_url,omitempty" validate:"omitempty"`
	Activity    string    `json:"activity" validate:"required,min=1"`
	Contacts    Contacts  `json:"contacts" validate:"required"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UpdateProfileRequest represents the request to update profile
type UpdateProfileRequest struct {
	Name        string   `json:"name" validate:"required,min=1,max=255"`
	Description string   `json:"description" validate:"required,min=1"`
	PhotoURL    string   `json:"photo_url,omitempty" validate:"omitempty"`
	Activity    string   `json:"activity" validate:"required,min=1"`
	Contacts    Contacts `json:"contacts" validate:"required"`
}
