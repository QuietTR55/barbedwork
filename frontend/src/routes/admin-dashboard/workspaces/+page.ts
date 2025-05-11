import { backendUrl } from '$lib/stores/backend';
import { authFetch } from '$lib/utilities/authFetch';
import { get } from 'svelte/store';
import type { PageLoad } from '../../$types';
import type { Workspace } from '$lib/models/workspace';

export const load: PageLoad = async ({ params }) => {
	const endPoint = get(backendUrl) + '/api/admin/workspaces';
	const response = await authFetch(endPoint, {
		method: 'GET'
	});

	if (!response.ok) {
		throw new Error('Failed to fetch workspaces');
	}

	const workspaces = (await response.json()) as Workspace[];
	return { workspaces: workspaces as Workspace[] };
};
