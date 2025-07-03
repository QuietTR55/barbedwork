CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    image_path VARCHAR(255)
);

CREATE INDEX IF NOT EXISTS idx_users_username ON users (username);

CREATE TABLE IF NOT EXISTS permissions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS user_permissions (
    user_id UUID NOT NULL REFERENCES users(id),
    permission_id INT NOT NULL REFERENCES permissions(id),
    workspace_id UUID NOT NULL REFERENCES workspaces(id),
    PRIMARY KEY (user_id, permission_id)
);

CREATE TABLE IF NOT EXISTS workspaces (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    image_path VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS workspace_users (
    workspace_id UUID NOT NULL REFERENCES workspaces(id),
    user_id UUID NOT NULL REFERENCES users(id),
    PRIMARY KEY (workspace_id, user_id)
);

CREATE TABLE IF NOT EXISTS teams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS workspace_teams (
    workspace_id UUID NOT NULL REFERENCES workspaces(id),
    team_id UUID NOT NULL REFERENCES teams(id),
    PRIMARY KEY (workspace_id, team_id)
);

CREATE TABLE IF NOT EXISTS team_users (
    team_id UUID NOT NULL REFERENCES teams(id),
    user_id UUID NOT NULL REFERENCES users(id),
    PRIMARY KEY (team_id, user_id)
);

INSERT INTO permissions (name) VALUES ('admin');
INSERT INTO permissions (name) VALUES ('user');

CREATE TABLE IF NOT EXISTS workspace_channels (
    id SERIAL PRIMARY KEY,
    workspace_id UUID NOT NULL REFERENCES workspaces(id),
    channel_name VARCHAR(255) NOT NULL,
    channel_emoji VARCHAR(255) NOT NULL DEFAULT '💬'
);

CREATE TABLE IF NOT EXISTS workspace_channel_messages (
    id SERIAL PRIMARY KEY,
    workspace_id UUID NOT NULL REFERENCES workspaces(id),
    channel_id INT NOT NULL REFERENCES workspace_channels(id),
    user_id UUID NOT NULL REFERENCES users(id),
    message TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS workspace_channel_message_reactions (
    id SERIAL PRIMARY KEY,
    message_id INT NOT NULL REFERENCES workspace_channel_messages(id),
    user_id UUID NOT NULL REFERENCES users(id),
    reaction VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS workspace_channel_message_replies (
    id SERIAL PRIMARY KEY,
    message_id INT NOT NULL REFERENCES workspace_channel_messages(id),
    user_id UUID NOT NULL REFERENCES users(id),
    reply TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);