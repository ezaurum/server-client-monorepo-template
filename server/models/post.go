package models

type Post struct {
	Model64int
	Title  string `json:"title"`
	UserID int64  `json:"userID"`
}
