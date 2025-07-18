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
    var channels []models.WorkspaceChannel = make([]models.WorkspaceChannel, 0)
    
    // Maps to track unique users and channels to avoid duplicates
    userMap := make(map[string]models.User)
    channelMap := make(map[string]models.WorkspaceChannel)

    query := `
        SELECT w.id, w.name, w.image_path,
               u.id, u.username, u.image_path,
               c.id, c.channel_name, c.workspace_id, c.channel_emoji
        FROM workspaces w
        LEFT JOIN workspace_users wu ON w.id = wu.workspace_id
        LEFT JOIN users u ON wu.user_id = u.id
        LEFT JOIN workspace_channels c ON w.id = c.workspace_id
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
        var tempWorkspaceImagePath sql.NullString

        // Use pgtype for nullable user fields, matching your pgx driver
        var userID pgtype.UUID
        var userName pgtype.Text
        var userImagePath pgtype.Text

        // Add channel variables to match the query
        var channelID sql.NullInt32
        var channelName sql.NullString
        var channelWorkspaceID pgtype.UUID
        var channelEmoji sql.NullString

        err := rows.Scan(
            &tempWorkspaceId,
            &tempWorkspaceName,
            &tempWorkspaceImagePath,
            &userID,
            &userName,
            &userImagePath,
            &channelID,          // Scan channel ID
            &channelName,        // Scan channel name
            &channelWorkspaceID, // Scan channel workspace ID
            &channelEmoji,       // Scan channel emoji
        )
        if err != nil {
            return nil, fmt.Errorf("scanning workspace and user row: %w", err)
        }

        if !firstRowProcessed {
            workspaceData.Id = tempWorkspaceId
            workspaceData.Name = tempWorkspaceName
            if tempWorkspaceImagePath.Valid {
                workspaceData.ImagePath = tempWorkspaceImagePath
            } else {
                workspaceData.ImagePath = sql.NullString{}
            }
            firstRowProcessed = true
        }

        // Only create and add user if userID is valid (i.e., a user was found by the LEFT JOIN)
        if userID.Valid {
            userUUID := userID.Bytes
            userKey := fmt.Sprintf("%x", userUUID)
            
            if _, exists := userMap[userKey]; !exists {
                currentUser := models.User{
                    Id: userID,
                }
                if userName.Valid {
                    currentUser.Username = userName.String
                }
                if userImagePath.Valid {
                    currentUser.ImagePath = userImagePath
                }
                userMap[userKey] = currentUser
            }
        }

        // Handle channel data if channel ID is valid
        if channelID.Valid {
            channelKey := fmt.Sprintf("%d", channelID.Int32)
            
            if _, exists := channelMap[channelKey]; !exists {
                currentChannel := models.WorkspaceChannel{
                    ID: channelKey,
                }
                if channelName.Valid {
                    currentChannel.Name = channelName.String
                }
                if channelEmoji.Valid {
                    currentChannel.Emoji = channelEmoji.String
                }
                channelMap[channelKey] = currentChannel
            }
        }
    }

    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("iteration error over workspace and user rows: %w", err)
    }

    if !firstRowProcessed {
        return nil, pgx.ErrNoRows // Workspace itself not found
    }

    // Convert maps to slices
    for _, user := range userMap {
        users = append(users, user)
    }
    for _, channel := range channelMap {
        channels = append(channels, channel)
    }

    workspaceData.Users = users
    workspaceData.Channels = channels
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

    userPermissionQuery := `
        INSERT INTO user_permissions (user_id, workspace_id, permission)
        VALUES ($1, $2, 'user')
    `
    _, err = repo.db.Exec(ctx, userPermissionQuery, userID, workspaceID)
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

func (repo *WorkspaceRepo) CreateChannel(ctx context.Context, workspaceId string, channelName string, channelEmoji string) (*models.WorkspaceChannel, error) {
    // Set default values if not provided
    if channelName == "" {
        channelName = "general"
    }
    if channelEmoji == "" {
        channelEmoji = "💬"
    }
    
    query := `
        INSERT INTO workspace_channels (workspace_id, channel_name, channel_emoji)
        VALUES ($1, $2, $3)
        RETURNING id, channel_name, channel_emoji
    `
    var channelId int32
    var returnedChannelName string
    var returnedChannelEmoji string
    err := repo.db.QueryRow(ctx, query, workspaceId, channelName, channelEmoji).Scan(&channelId, &returnedChannelName, &returnedChannelEmoji)
    if err != nil {
        return nil, err
    }

    return &models.WorkspaceChannel{
        ID:           fmt.Sprintf("%d", channelId),
        Name:         returnedChannelName,
        Emoji:        returnedChannelEmoji,
    }, nil
}   

func (repo *WorkspaceRepo) GetAvailablePermissions(ctx context.Context) ([]string, error) {
    query := `
        SELECT name
        FROM permissions
    `
    rows, err := repo.db.Query(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var permissions []string
    for rows.Next() {
        var permission string
        err := rows.Scan(&permission)
        if err != nil {
            return nil, err // Return any error encountered during scanning
        }
        permissions = append(permissions, permission)
    }

    if err = rows.Err(); err != nil {
        return nil, err // Check for errors during iteration
    }
    return permissions, nil
}