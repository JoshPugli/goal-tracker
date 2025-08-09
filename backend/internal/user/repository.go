package user

import (
	"database/sql"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(email, firstName, password string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	query := `
		INSERT INTO users (email, first_name, password) 
		VALUES ($1, $2, $3) 
		RETURNING id`

	var id string
	err = r.db.QueryRow(query, email, firstName, string(hashedPassword)).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &User{
		ID:        id,
		Email:     email,
		FirstName: firstName,
	}, nil
}

func (r *Repository) GetUserByEmail(email string) (*User, error) {
	query := `SELECT id, email, first_name, password FROM users WHERE email = $1`

	user := &User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.FirstName, &user.Password,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *Repository) GetUserByID(id string) (*User, error) {
	query := `SELECT id, email, first_name, last_name FROM users WHERE id = $1`

	user := &User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.FirstName,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *Repository) ValidatePassword(email, password string) (*User, error) {
	user, err := r.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	return user, nil
}