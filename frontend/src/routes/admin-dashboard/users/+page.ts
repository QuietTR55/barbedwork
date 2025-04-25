import { backendUrl } from '$lib/stores/backend';
import { get } from 'svelte/store';
import { getAccessToken } from '$lib/stores/authentication';
import type { PageLoad } from './$types';

const apiUrl = get(backendUrl);

export const load = (async ({ fetch }) => {}) satisfies PageLoad;
