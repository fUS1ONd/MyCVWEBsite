// Package domain contains core business entities and types
package domain

// Profile represents user profile information displayed on CV page
type Profile struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Skills      []string `json:"skills"`
	Contacts    Contacts `json:"contacts"`
}

// Contacts represents contact information
type Contacts struct {
	Email    string `json:"email"`
	Github   string `json:"github"`
	LinkedIn string `json:"linkedin"`
}
