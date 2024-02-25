package models

type User struct {
	ID   string `form:"id" json:"id"`
	Name string `form:"name" json:"name"`
}
