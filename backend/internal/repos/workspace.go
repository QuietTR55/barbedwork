package repos

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkspaceRepo struct {
	db *pgxpool.Pool
}

func NewWorkspaceRepo(db *pgxpool.Pool) *WorkspaceRepo {
	return &WorkspaceRepo{db: db}
}

func (repo *WorkspaceRepo) CreateWorkspace(ctx context.Context, workspaceName string, workspaceImagePath string) (*models.Workspace, error) {
	var workspaceID uuid.UUID
	var createdWorkspace models.Workspace
	fmt.Println("Creating workspace with name:", workspaceName, "and image path:", workspaceImagePath)
	query := `
		INSERT INTO workspaces (name, image_path)
		VALUES ($1, $2)
		RETURNING id, name, image_path
	`
	err := repo.db.QueryRow(ctx, query, workspaceName, workspaceImagePath).Scan(&workspaceID, &createdWorkspace.Name, &createdWorkspace.ImagePath)
	if err != nil {
		return nil, err
	}
	return &createdWorkspace, nil
}

func (repo *WorkspaceRepo) GetWorkspace(ctx context.Context, workspaceID string) (*models.WorkspaceFullData, error) {
    var workspaceData models.WorkspaceFullData
    var users []models.User = make([]models.User, 0)

    query := `
        SELECT w.id, w.name, w.image_path,
               u.id, u.username, u.image_path
        FROM workspaces w
        LEFT JOIN workspace_users wu ON w.id = wu.workspace_id
        LEFT JOIN users u ON wu.user_id = u.id
        WHERE w.id = $1
    `

    rows, err := repo.db.Query(ctx, query, workspaceID)
    if err != nil {
        return nil, fmt.Errorf("querying workspace and users: %w", err)
    }
    defer rows.Close()

    firstRowProcessed := false
    for rows.Next() {
        var tempWorkspaceId uuid.UUID
        var tempWorkspaceName string
        var tempWorkspaceImagePath sql.NullString // This was correct

        // Use pgtype for nullable user fields, matching your pgx driver
        var userID pgtype.UUID
        var userName pgtype.Text
        var userImagePath pgtype.Text

        err := rows.Scan(
            &tempWorkspaceId,
            &tempWorkspaceName,
            &tempWorkspaceImagePath, // Scan workspace image path
            &userID,                 // Scan user ID into nullable type
            &userName,               // Scan username into nullable type
            &userImagePath,          // Scan user image path into nullable type
        )
        if err != nil {
            return nil, fmt.Errorf("scanning workspace and user row: %w", err)
        }

        if !firstRowProcessed {
            workspaceData.Id = tempWorkspaceId
            workspaceData.Name = tempWorkspaceName
            if tempWorkspaceImagePath.Valid { // Check if workspace image path is not null
                workspaceData.ImagePath = tempWorkspaceImagePath
            } else {
                workspaceData.ImagePath = sql.NullString{}
            }
            firstRowProcessed = true
        }

        // Only create and add user if userID is valid (i.e., a user was found by the LEFT JOIN)
        if userID.Valid {
            currentUser := models.User{
                Id: userID, // Get the actual UUID value
            }
            if userName.Valid {
                currentUser.Username = userName.String
            }
            if userImagePath.Valid {
                currentUser.ImagePath = userImagePath
            }
            users = append(users, currentUser)
        }
    }

    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("iteration error over workspace and user rows: %w", err)
    }

    if !firstRowProcessed {
        return nil, pgx.ErrNoRows // Workspace itself not found
    }

    workspaceData.Users = users
    return &workspaceData, nil
}

func (repo *WorkspaceRepo) GetAllWorkspaces(ctx context.Context) ([]*models.Workspace, error) {
	query := `
		SELECT id, name, image_path
		FROM workspaces
	`
	rows, err := repo.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workspaces []*models.Workspace
	for rows.Next() {
		var workspace models.Workspace
		err := rows.Scan(&workspace.Id, &workspace.Name, &workspace.ImagePath)
		if err != nil {
			return nil, err // Return any error encountered during scanning
		}
		workspaces = append(workspaces, &workspace)
	}

	if err = rows.Err(); err != nil {
		return nil, err // Check for errors during iteration
	}
	return workspaces, nil
}

func (repo *WorkspaceRepo) AddUserToWorkspace(ctx context.Context, userID string, workspaceID string) error {
    query := `
        INSERT INTO workspace_users (user_id, workspace_id)
        VALUES ($1, $2)
        ON CONFLICT (user_id, workspace_id) DO NOTHING
    `
    _, err := repo.db.Exec(ctx, query, userID, workspaceID)
    if err != nil {
        return err
    }
    return nil
}

func (repo *WorkspaceRepo) GetUserWorkspaces(ctx context.Context, userId string) ([]*models.Workspace, error) {
    query := `
        SELECT w.id, w.name, w.image_path
        FROM workspaces w
        JOIN workspace_users wu ON w.id = wu.workspace_id
        WHERE wu.user_id = $1
    `
    rows, err := repo.db.Query(ctx, query, userId)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var workspaces []*models.Workspace
    for rows.Next() {
        var workspace models.Workspace
        err := rows.Scan(&workspace.Id, &workspace.Name, &workspace.ImagePath)
        if err != nil {
            return nil, err // Return any error encountered during scanning
        }
        workspaces = append(workspaces, &workspace)
    }

    if err = rows.Err(); err != nil {
        return nil, err // Check for errors during iteration
    }
    return workspaces, nil
}