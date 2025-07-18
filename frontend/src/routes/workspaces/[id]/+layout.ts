import type { Workspace } from '$lib/models/workspace';
import { authFetch } from '$lib/utilities/authFetch';
import { get } from 'svelte/store';
import type { LayoutLoad } from './$types';
import { backendUrl } from '$lib/stores/backend';
import { workspaceActions } from '$lib/stores/workspace';


export const load = (async ({ params }) => {
    const workspaceId: string = params.id;
    if(!workspaceId) {
        throw new Error('Workspace ID is required');
    }
    const backend = get(backendUrl);
    const workspaceFetch = await authFetch(`${backend}/api/workspaces/${workspaceId}`)
    if (!workspaceFetch.ok) {
        throw new Error(`Failed to fetch workspace with ID ${workspaceId}`);
    }
    const workspace: Workspace = await workspaceFetch.json();
    workspaceActions.setWorkspace(workspace);
    console.log('Loaded workspace:', workspace);
    return {
        workspaceId,
        workspace
    };
}) satisfies LayoutLoad;