package user

type User struct {
	ID        string `json:"id" db:"id"`
	Email     string `json:"email" db:"email"`
	Password  string `json:"-" db:"password"` // Never expose in JSON
	FirstName string `json:"first_name" db:"first_name"`
}