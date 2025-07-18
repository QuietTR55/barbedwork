import type { Role, Permission, WorkspaceUserRole } from '$lib/models/role';
import type { User } from '$lib/models/user';
import { authFetch } from '$lib/utilities/authFetch';
import { get } from 'svelte/store';
import { backendUrl } from '$lib/stores/backend';

const apiUrl = get(backendUrl) + '/api';

// Permissions from your database schema
export const availablePermissions: Permission[] = [
  { id: 1, name: 'workspace:manage-users', description: 'Manage users in the workspace (invite, remove, edit user info)' },
  { id: 2, name: 'workspace:manage-teams', description: 'Create, edit, and delete teams' },
  { id: 3, name: 'workspace:manage-channels', description: 'Create, edit, and delete channels' },
  { id: 4, name: 'workspace:manage-roles', description: 'Create, edit, delete roles and assign permissions to roles' },
  { id: 5, name: 'workspace:send-messages', description: 'Send messages in channels' },
  { id: 6, name: 'workspace:delete-any-message', description: 'Delete any message in channels' },
  { id: 7, name: 'workspace:delete-own-message', description: 'Delete own messages' },
  { id: 8, name: 'workspace:edit-any-message', description: 'Edit any message in channels' },
  { id: 9, name: 'workspace:edit-own-message', description: 'Edit own messages' },
  { id: 10, name: 'workspace:pin-messages', description: 'Pin and unpin messages' },
  { id: 11, name: 'workspace:manage-workspace', description: 'Edit workspace settings, name, image, etc.' },
  { id: 12, name: 'workspace:view-channels', description: 'View and read messages in channels' },
  { id: 13, name: 'workspace:create-channels', description: 'Create new channels' },
  { id: 14, name: 'workspace:archive-channels', description: 'Archive and unarchive channels' },
  { id: 15, name: 'workspace:manage-reactions', description: 'Add and remove reactions to messages' },
  { id: 16, name: 'workspace:upload-files', description: 'Upload files to channels' }
];

export class RoleService {
  static async getPermissions(): Promise<Permission[]> {
    try {
      const response = await authFetch(`${apiUrl}/permissions`);
      if (!response.ok) throw new Error('Failed to fetch permissions');
      return await response.json();
    } catch (error) {
      console.error('Failed to fetch permissions:', error);
      return availablePermissions; // Fallback to static data
    }
  }

  static async getRoles(): Promise<Role[]> {
    try {
      const response = await authFetch(`${apiUrl}/roles`);
      if (!response.ok) throw new Error('Failed to fetch roles');
      return await response.json();
    } catch (error) {
      console.error('Failed to fetch roles:', error);
      return []
    }
  }

  static async getWorkspaceUserRoles(workspaceId: string): Promise<WorkspaceUserRole[]> {
    try {
      const response = await authFetch(`${apiUrl}/workspaces/${workspaceId}/user-roles`);
      if (!response.ok) throw new Error('Failed to fetch workspace user roles');
      return await response.json();
    } catch (error) {
      console.error('Failed to fetch workspace user roles:', error);
      return []
    }
  }

  static async createRole(role: Omit<Role, 'id' | 'created_at' | 'updated_at'>): Promise<Role> {
    try {
      const payload = {
        name: role.name,
        description: role.description,
        permission_ids: role.permissions.map(p => p.id)
      };

      const response = await authFetch(`${apiUrl}/roles`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload)
      });

      if (!response.ok) {
        const error = await response.text();
        throw new Error(error || 'Failed to create role');
      }

      return await response.json();
    } catch (error) {
      console.error('Failed to create role:', error);
      throw error;
    }
  }

  static async updateRole(roleId: number, updates: Partial<Omit<Role, 'id' | 'created_at'>>): Promise<Role> {
    try {
      const payload = {
        name: updates.name,
        description: updates.description,
        permission_ids: updates.permissions?.map(p => p.id) || []
      };

      const response = await authFetch(`${apiUrl}/roles/${roleId}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload)
      });

      if (!response.ok) {
        const error = await response.text();
        throw new Error(error || 'Failed to update role');
      }

      return await response.json();
    } catch (error) {
      console.error('Failed to update role:', error);
      throw error;
    }
  }

  static async deleteRole(roleId: number): Promise<void> {
    try {
      const response = await authFetch(`${apiUrl}/roles/${roleId}`, {
        method: 'DELETE'
      });

      if (!response.ok) {
        const error = await response.text();
        throw new Error(error || 'Failed to delete role');
      }
    } catch (error) {
      console.error('Failed to delete role:', error);
      throw error;
    }
  }

  static async assignRoleToUser(workspaceId: string, userId: string, roleId: number): Promise<WorkspaceUserRole> {
    try {
      const payload = {
        user_id: userId,
        role_id: roleId
      };

      const response = await authFetch(`${apiUrl}/workspaces/${workspaceId}/user-roles`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(payload)
      });

      if (!response.ok) {
        const error = await response.text();
        throw new Error(error || 'Failed to assign role');
      }

      return await response.json();
    } catch (error) {
      console.error('Failed to assign role:', error);
      throw error;
    }
  }

  static async removeRoleFromUser(workspaceId: string, userId: string, roleId: number): Promise<void> {
    try {
      const response = await authFetch(
        `${apiUrl}/workspaces/${workspaceId}/user-roles?user_id=${userId}&role_id=${roleId}`, 
        {
          method: 'DELETE'
        }
      );

      if (!response.ok) {
        const error = await response.text();
        throw new Error(error || 'Failed to remove role assignment');
      }
    } catch (error) {
      console.error('Failed to remove role assignment:', error);
      throw error;
    }
  }
}