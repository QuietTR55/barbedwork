import type { Channel } from '$lib/models/channel';
import type { User } from '$lib/models/user';
import type { Workspace } from '$lib/models/workspace';
import { backendUrl } from '$lib/stores/backend';
import { authFetch } from '$lib/utilities/authFetch';
import { get } from 'svelte/store';

export async function addUserToWorkspace(workspaceId: string, userId: string): Promise<boolean> {
	const endPoint = get(backendUrl) + '/api/admin/workspaces/' + workspaceId + '/users/' + userId;
	const response = await authFetch(endPoint, {
		method: 'POST'
	});

	if (!response.ok) {
		throw new Error('Failed to add user to workspace');
	}
	return true;
}

export async function getUserWorkspaces(): Promise<Workspace[]> {
	const endPoint = get(backendUrl) + '/api/workspaces';
	const response = await authFetch(endPoint, {
		method: 'GET'
	});

	if (!response.ok) {
		throw new Error('Failed to fetch workspaces');
	}

	const data = await response.json();
	return data;
}

export async function createChannel(workspaceId: string, channelData?: { name: string; emoji: string }): Promise<Channel> {
	const endPoint = get(backendUrl) + '/api/workspaces/' + workspaceId + '/channels';
	const response = await authFetch(endPoint, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify({
			name: channelData?.name || 'New Channel',
			emoji: channelData?.emoji || '🔔'
		})
	});

	if (!response.ok) {
		throw new Error('Failed to create channel');
	}

	const data = await response.json();
	return data;
}