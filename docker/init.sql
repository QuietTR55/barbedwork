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
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT
);

CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS role_permissions (
    role_id INT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id INT NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE IF NOT EXISTS workspace_user_roles (
    workspace_id UUID NOT NULL REFERENCES workspaces(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id INT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (workspace_id, user_id, role_id)
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

-- Insert comprehensive permissions
INSERT INTO permissions (name, description) VALUES 
    ('workspace:manage-users', 'Manage users in the workspace (invite, remove, edit user info)') ON CONFLICT (name) DO NOTHING;
INSERT INTO permissions (name, description) VALUES 
    ('workspace:manage-teams', 'Create, edit, and delete teams') ON CONFLICT (name) DO NOTHING;
INSERT INTO permissions (name, description) VALUES 
    ('workspace:manage-channels', 'Create, edit, and delete channels') ON CONFLICT (name) DO NOTHING;
INSERT INTO permissions (name, description) VALUES 
    ('workspace:manage-roles', 'Create, edit, delete roles and assign permissions to roles') ON CONFLICT (name) DO NOTHING;
INSERT INTO permissions (name, description) VALUES 
    ('workspace:send-messages', 'Send messages in channels') ON CONFLICT (name) DO NOTHING;
INSERT INTO permissions (name, description) VALUES 
    ('workspace:delete-any-message', 'Delete any message in channels') ON CONFLICT (name) DO NOTHING;
INSERT INTO permissions (name, description) VALUES 
    ('workspace:delete-own-message', 'Delete own messages') ON CONFLICT (name) DO NOTHING;
INSERT INTO permissions (name, description) VALUES 
    ('workspace:edit-any-message', 'Edit any message in channels') ON CONFLICT (name) DO NOTHING;
INSERT INTO permissions (name, description) VALUES 
    ('workspace:edit-own-message', 'Edit own messages') ON CONFLICT (name) DO NOTHING;
INSERT INTO permissions (name, description) VALUES 
    ('workspace:pin-messages', 'Pin and unpin messages') ON CONFLICT (name) DO NOTHING;
INSERT INTO permissions (name, description) VALUES 
    ('workspace:manage-workspace', 'Edit workspace settings, name, image, etc.') ON CONFLICT (name) DO NOTHING;
INSERT INTO permissions (name, description) VALUES 
    ('workspace:view-channels', 'View and read messages in channels') ON CONFLICT (name) DO NOTHING;
INSERT INTO permissions (name, description) VALUES 
    ('workspace:create-channels', 'Create new channels') ON CONFLICT (name) DO NOTHING;
INSERT INTO permissions (name, description) VALUES 
    ('workspace:archive-channels', 'Archive and unarchive channels') ON CONFLICT (name) DO NOTHING;
INSERT INTO permissions (name, description) VALUES 
    ('workspace:manage-reactions', 'Add and remove reactions to messages') ON CONFLICT (name) DO NOTHING;
INSERT INTO permissions (name, description) VALUES 
    ('workspace:upload-files', 'Upload files to channels') ON CONFLICT (name) DO NOTHING;

-- Insert default roles
INSERT INTO roles (name, description) VALUES 
    ('Owner', 'Full access to all workspace features and settings') ON CONFLICT (name) DO NOTHING;
INSERT INTO roles (name, description) VALUES 
    ('Admin', 'Administrative access with user and team management capabilities') ON CONFLICT (name) DO NOTHING;
INSERT INTO roles (name, description) VALUES 
    ('Moderator', 'Moderate messages and manage channels') ON CONFLICT (name) DO NOTHING;
INSERT INTO roles (name, description) VALUES 
    ('Member', 'Standard user with basic messaging capabilities') ON CONFLICT (name) DO NOTHING;
INSERT INTO roles (name, description) VALUES 
    ('Guest', 'Limited access for temporary users') ON CONFLICT (name) DO NOTHING;

-- Assign permissions to Owner role (all permissions)
INSERT INTO role_permissions (role_id, permission_id) 
SELECT r.id, p.id FROM roles r, permissions p 
WHERE r.name = 'Owner' ON CONFLICT DO NOTHING;

-- Assign permissions to Admin role
INSERT INTO role_permissions (role_id, permission_id) 
SELECT r.id, p.id FROM roles r, permissions p 
WHERE r.name = 'Admin' AND p.name IN (
    'workspace:manage-users',
    'workspace:manage-teams', 
    'workspace:manage-channels',
    'workspace:send-messages',
    'workspace:delete-any-message',
    'workspace:delete-own-message',
    'workspace:edit-any-message',
    'workspace:edit-own-message',
    'workspace:pin-messages',
    'workspace:view-channels',
    'workspace:create-channels',
    'workspace:archive-channels',
    'workspace:manage-reactions',
    'workspace:upload-files'
) ON CONFLICT DO NOTHING;

-- Assign permissions to Moderator role
INSERT INTO role_permissions (role_id, permission_id) 
SELECT r.id, p.id FROM roles r, permissions p 
WHERE r.name = 'Moderator' AND p.name IN (
    'workspace:manage-channels',
    'workspace:send-messages',
    'workspace:delete-any-message',
    'workspace:delete-own-message',
    'workspace:edit-any-message',
    'workspace:edit-own-message',
    'workspace:pin-messages',
    'workspace:view-channels',
    'workspace:create-channels',
    'workspace:archive-channels',
    'workspace:manage-reactions',
    'workspace:upload-files'
) ON CONFLICT DO NOTHING;

-- Assign permissions to Member role
INSERT INTO role_permissions (role_id, permission_id) 
SELECT r.id, p.id FROM roles r, permissions p 
WHERE r.name = 'Member' AND p.name IN (
    'workspace:send-messages',
    'workspace:delete-own-message',
    'workspace:edit-own-message',
    'workspace:view-channels',
    'workspace:manage-reactions',
    'workspace:upload-files'
) ON CONFLICT DO NOTHING;

-- Assign permissions to Guest role
INSERT INTO role_permissions (role_id, permission_id) 
SELECT r.id, p.id FROM roles r, permissions p 
WHERE r.name = 'Guest' AND p.name IN (
    'workspace:send-messages',
    'workspace:delete-own-message',
    'workspace:edit-own-message',
    'workspace:view-channels',
    'workspace:manage-reactions'
) ON CONFLICT DO NOTHING;

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

CREATE INDEX IF NOT EXISTS idx_workspace_channel_messages_created_at ON workspace_channel_messages (created_at);

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