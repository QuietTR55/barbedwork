import { writable } from 'svelte/store';
import type { Workspace } from '$lib/models/workspace';
import type { Channel } from '$lib/models/channel';
import type { User } from '$lib/models/user';

// Create a writable store for the active workspace
export const activeWorkspace = writable<Workspace | null>(null);

// Helper functions to update specific parts of the workspace
export const workspaceActions = {
    // Set the entire workspace
    setWorkspace: (workspace: Workspace) => {
        activeWorkspace.set(workspace);
    },

    // Clear the workspace
    clearWorkspace: () => {
        activeWorkspace.set(null);
    },

    // Update workspace name
    updateWorkspaceName: (name: string) => {
        activeWorkspace.update(workspace => {
            if (workspace) {
                return { ...workspace, Name: name };
            }
            return workspace;
        });
    },

    // Update workspace image
    updateWorkspaceImage: (imagePath: string, isValid: boolean = true) => {
        activeWorkspace.update(workspace => {
            if (workspace) {
                return {
                    ...workspace,
                    ImagePath: { String: imagePath, Valid: isValid }
                };
            }
            return workspace;
        });
    },

    // Add a new channel
    addChannel: (channel: Channel) => {
        activeWorkspace.update(workspace => {
            if (workspace) {
                return {
                    ...workspace,
                    Channels: [...workspace.Channels, channel]
                };
            }
            return workspace;
        });
    },

    // Remove a channel
    removeChannel: (channelId: string) => {
        activeWorkspace.update(workspace => {
            if (workspace) {
                return {
                    ...workspace,
                    Channels: workspace.Channels.filter(channel => channel.Id !== channelId)
                };
            }
            return workspace;
        });
    },

    // Update a channel
    updateChannel: (channelId: string, updates: Partial<Channel>) => {
        activeWorkspace.update(workspace => {
            if (workspace) {
                return {
                    ...workspace,
                    Channels: workspace.Channels.map(channel =>
                        channel.Id === channelId ? { ...channel, ...updates } : channel
                    )
                };
            }
            return workspace;
        });
    },

    // Add a user to workspace
    addUser: (user: User) => {
        activeWorkspace.update(workspace => {
            if (workspace) {
                // Check if user already exists to avoid duplicates
                const userExists = workspace.Users.some(existingUser => existingUser.id === user.id);
                if (!userExists) {
                    return {
                        ...workspace,
                        Users: [...workspace.Users, user]
                    };
                }
            }
            return workspace;
        });
    },

    // Remove a user from workspace
    removeUser: (userId: string) => {
        activeWorkspace.update(workspace => {
            if (workspace) {
                return {
                    ...workspace,
                    Users: workspace.Users.filter(user => user.id !== userId)
                };
            }
            return workspace;
        });
    },

    // Update a user
    updateUser: (userId: string, updates: Partial<User>) => {
        activeWorkspace.update(workspace => {
            if (workspace) {
                return {
                    ...workspace,
                    Users: workspace.Users.map(user =>
                        user.id === userId ? { ...user, ...updates } : user
                    )
                };
            }
            return workspace;
        });
    }
};

// Derived stores for specific workspace data
export const workspaceChannels = writable<Channel[]>([]);
export const workspaceUsers = writable<User[]>([]);
export const workspaceName = writable<string>('');

// Subscribe to workspace changes and update derived stores
activeWorkspace.subscribe(workspace => {
    if (workspace) {
        workspaceChannels.set(workspace.Channels);
        workspaceUsers.set(workspace.Users);
        workspaceName.set(workspace.Name);
    } else {
        workspaceChannels.set([]);
        workspaceUsers.set([]);
        workspaceName.set('');
    }
});