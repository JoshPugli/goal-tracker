package models

type User struct {
	ID        string `json:"id" db:"id"`
	Email     string `json:"email" db:"email"`
	Username  string `json:"username" db:"username"`
	Password  string `json:"-" db:"password"` // Never expose in JSON
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
}