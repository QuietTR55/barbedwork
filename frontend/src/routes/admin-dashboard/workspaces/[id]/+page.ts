import { get } from 'svelte/store';
import type { PageLoad } from '../$types';
import { backendUrl } from '$lib/stores/backend';
import { authFetch } from '$lib/utilities/authFetch';
import { getAllUsers } from '$lib/services/adminDashboard';
import type { User } from '$lib/models/user';
import type { Workspace } from '$lib/models/workspace'; // Import Workspace type

export const load: PageLoad = async ({ params }) => {
	try {
		const workspaceId = params.id;
		const endPoint = get(backendUrl) + '/api/admin/workspaces/' + workspaceId;
		const response = await authFetch(endPoint, {
			method: 'GET'
		});

		if (!response.ok) {
			console.error(
				'Failed to fetch workspace:',
				response.status,
				await response.text().catch(() => '')
			);
			// Return an empty workspace or an error structure
			return {
				workspace: { Id: workspaceId, Name: 'Error Loading Workspace', users: [] },
				allUsers: []
			};
		}

		const backendWorkspaceData = await response.json();

		const workspace: Workspace = {
			Id: backendWorkspaceData.Id,
			Name: backendWorkspaceData.Name,
			ImagePath: backendWorkspaceData.ImagePath,
			Users: backendWorkspaceData.Users || []
		};

		const { users: fetchedAllUsers, error: fetchAllUsersError } = await getAllUsers();
		const allUsers = fetchedAllUsers || [];

		if (fetchAllUsersError) {
			console.error('Error fetching all users:', fetchAllUsersError);
			return { workspace, allUsers: [] };
		}

		return { workspace, allUsers };
	} catch (error) {
		console.error('Error loading workspace data in +page.ts:', error);
		return {
			workspace: { Id: params.id, Name: 'Error Loading Workspace', users: [] },
			allUsers: [],
			error: 'Failed to load workspace details'
		};
	}
};
