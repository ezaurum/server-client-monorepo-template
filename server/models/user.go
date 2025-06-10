package models

// User model
type User struct {
	Model64int
	Name  string `json:"name"`
	Email string `json:"email"`
	Posts []Post `json:"posts"`
}
