import { get } from 'svelte/store';
import type { PageLoad } from '../$types';
import { backendUrl } from '$lib/stores/backend';
import { authFetch } from '$lib/utilities/authFetch';

export const load: PageLoad = async ({ params }) => {
	try {
		const workspaceId = params.id;
		const endPoint = get(backendUrl) + '/api/admin/workspaces/' + workspaceId;
		const response = await authFetch(endPoint, {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json'
			},
			credentials: 'include'
		});

		const workspace = await response.json();
		return { workspace };
	} catch (error) {
		console.error('Error loading workspace:', error);
		return {};
	}
};
