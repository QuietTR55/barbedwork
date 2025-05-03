package repos

import (
	"backend/internal/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) CreateUser(ctx context.Context, username string, passwordHash string) (*models.User, error) {
	var query string = `
	INSERT INTO users (username, password_hash)
	VALUES ($1, $2)
	RETURNING id, username
	`

	var id string
	var returnedUsername string
	err := r.db.QueryRow(ctx, query, username, passwordHash).Scan(&id, &returnedUsername)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:       id,
		Username: returnedUsername,
	}

	return user, nil
}

func (r *UserRepo) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	var query string = `
	SELECT id, username FROM users
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*models.User{}
	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.ID, &user.Username)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

func (r *UserRepo) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	fmt.Println("GetUserByUsername called with username:", username)
	var query string = `
	SELECT id, username, password_hash FROM users WHERE username = $1
	`

	var id string
	var returnedUsername string
	var passwordHash string
	err := r.db.QueryRow(ctx, query, username).Scan(&id, &returnedUsername, &passwordHash)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:           id,
		Username:     returnedUsername,
		PasswordHash: passwordHash,
	}

	return user, nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	var query string = `
	SELECT id, username FROM users WHERE id = $1
	`

	var id string
	var returnedUsername string
	err := r.db.QueryRow(ctx, query, userID).Scan(&id, &returnedUsername)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:       id,
		Username: returnedUsername,
	}

	return user, nil
}
