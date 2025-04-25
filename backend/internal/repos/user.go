package repos

import (
	"backend/internal/models"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(username string, passwordHash string) error {	
	var query string = `
	INSERT INTO users (username, password_hash)
	VALUES ($1, $2)
	RETURNING id
	`

	var id string
	err := r.db.QueryRow(context.Background(), query, username, passwordHash).Scan(&id)
	if err != nil {
		return err
	}
	
	return nil
}

func (r *UserRepo) GetAllUsers() ([]*models.User, error) {
	var query string = `
	SELECT id, username, password_hash FROM users
	`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*models.User{}
	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.ID, &user.Username, &user.PasswordHash)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}
