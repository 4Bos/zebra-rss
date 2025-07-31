package users

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        int64
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Repository interface {
	CreateUser(email string, password string) (*User, error)
	VerifyCredentials(email string, password string) error
	GetUserById(id int64) (*User, error)
	GetUserByEmail(email string) (*User, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(email string, password string) (*User, error) {
	var userId int64

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, err
	}

	query := "INSERT INTO zebra.users (email, password) VALUES ($1, $2) RETURNING id"
	row := r.db.QueryRow(query, email, hashedPassword)
	err = row.Scan(&userId)

	if err != nil {
		return nil, err
	}

	return r.GetUserById(userId)
}

func (r *repository) VerifyCredentials(email string, password string) error {
	var hashedPassword string

	query := "SELECT password FROM zebra.users WHERE email = $1"
	row := r.db.QueryRow(query, email)
	err := row.Scan(&hashedPassword)

	if err != nil {
		return err
	}

	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (r *repository) GetUserById(id int64) (*User, error) {
	var user User

	query := "SELECT id, email, created_at, updated_at FROM zebra.users WHERE id = $1"
	row := r.db.QueryRow(query, id)
	err := row.Scan(&user.Id, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repository) GetUserByEmail(email string) (*User, error) {
	var user User

	query := "SELECT id, email, created_at, updated_at FROM zebra.users WHERE email = $1"
	row := r.db.QueryRow(query, email)
	err := row.Scan(&user.Id, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
