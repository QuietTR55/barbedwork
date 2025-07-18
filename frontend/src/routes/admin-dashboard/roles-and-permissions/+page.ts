import type { PageLoad } from './$types';
import { RoleService } from '$lib/services/roleService';
import type { Role, Permission, WorkspaceUserRole } from '$lib/models/role';

export const load = (async ({ params }) => {
    try {
        // Admin dashboard focuses on global roles and permissions management
        // Workspace-specific user role assignments should be managed within individual workspaces
        const [permissions, roles] = await Promise.all([
            RoleService.getPermissions(),
            RoleService.getRoles()
        ]);

        return {
            roles,
            workspaceUserRoles: [] as WorkspaceUserRole[], // Empty for admin dashboard
            permissions,
            workspaceId: null // No specific workspace context in admin dashboard
        };
    } catch (error) {
        console.error('Failed to load roles and permissions:', error);
        return {
            roles: [] as Role[],
            workspaceUserRoles: [] as WorkspaceUserRole[],
            permissions: [] as Permission[],
            workspaceId: null
        };
    }
}) satisfies PageLoad;