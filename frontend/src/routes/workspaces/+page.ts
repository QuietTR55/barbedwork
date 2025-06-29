import { getUserWorkspaces } from '$lib/services/workspace';
import type { PageLoad } from '../$types';

export const load: PageLoad = async ({ params }) => {
	let workspaces = await getUserWorkspaces();
	if (!workspaces) {
		workspaces = [];
	}
	return {
		workspaces
	};
};
