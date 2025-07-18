package repos

import (
	"backend/internal/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
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

	var id pgtype.UUID
	var returnedUsername string
	err := r.db.QueryRow(ctx, query, username, passwordHash).Scan(&id, &returnedUsername)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Id:       id,
		Username: returnedUsername,
	}

	return user, nil
}

func (r *UserRepo) GetAllUsers(ctx context.Context) ([]*models.User, error) {
	var query string = `
	SELECT DISTINCT u.id, u.username, u.image_path, 
	       COALESCE(ARRAY_AGG(p.name) FILTER (WHERE p.name IS NOT NULL), '{}') as permissions
	FROM users u
	LEFT JOIN user_permissions up ON u.id = up.user_id
	LEFT JOIN permissions p ON up.permission_id = p.id
	GROUP BY u.id, u.username, u.image_path
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []*models.User{}
	for rows.Next() {
		var user models.User
		var permissions []string
		err = rows.Scan(&user.Id, &user.Username, &user.ImagePath, &permissions)
		if err != nil {
			return nil, err
		}
		
		user.Permissions = permissions
		users = append(users, &user)
	}

	return users, nil
}

func (r *UserRepo) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	fmt.Println("GetUserByUsername called with username:", username)
	var query string = `
	SELECT u.id, u.username, u.password_hash, u.image_path,
	       COALESCE(ARRAY_AGG(p.name) FILTER (WHERE p.name IS NOT NULL), '{}') as permissions
	FROM users u
	LEFT JOIN user_permissions up ON u.id = up.user_id
	LEFT JOIN permissions p ON up.permission_id = p.id
	WHERE u.username = $1
	GROUP BY u.id, u.username, u.password_hash, u.image_path
	`

	var id pgtype.UUID
	var returnedUsername string
	var passwordHash string
	var imagePath pgtype.Text
	var permissions []string
	
	err := r.db.QueryRow(ctx, query, username).Scan(&id, &returnedUsername, &passwordHash, &imagePath, &permissions)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Id:           id,
		Username:     returnedUsername,
		PasswordHash: passwordHash,
		ImagePath:    imagePath,
		Permissions:  permissions,
	}

	return user, nil
}

func (r *UserRepo) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	var userQuery string = `
	SELECT u.id, u.username, u.image_path,
	       COALESCE(ARRAY_AGG(p.name) FILTER (WHERE p.name IS NOT NULL), '{}') as permissions
	FROM users u
	LEFT JOIN user_permissions up ON u.id = up.user_id
	LEFT JOIN permissions p ON up.permission_id = p.id
	WHERE u.id = $1
	GROUP BY u.id, u.username, u.image_path
	`

	var id pgtype.UUID
	var returnedUsername string
	var imagePath pgtype.Text
	var permissions []string
	
	err := r.db.QueryRow(ctx, userQuery, userID).Scan(&id, &returnedUsername, &imagePath, &permissions)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Id:          id,
		Username:    returnedUsername,
		ImagePath:   imagePath,
		Permissions: permissions,
	}

	return user, nil
}

func (r *UserRepo) GetUserPermissions(ctx context.Context, userID string, workspaceId string) ([]string, error) {
	var query string = `
	SELECT p.name 
	FROM user_permissions up
	JOIN permissions p ON up.permission_id = p.id
	WHERE up.user_id = $1 AND up.workspace_id = $2
	`

	rows, err := r.db.Query(ctx, query, userID, workspaceId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []string
	for rows.Next() {
		var permission string
		err = rows.Scan(&permission)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	return permissions, nil
}

func (r *UserRepo) GetAllUserPermissions(ctx context.Context, userID string) ([]string, error) {
	var query string = `
	SELECT DISTINCT p.name 
	FROM user_permissions up
	JOIN permissions p ON up.permission_id = p.id
	WHERE up.user_id = $1
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []string
	for rows.Next() {
		var permission string
		err = rows.Scan(&permission)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	return permissions, nil
}