import { backendUrl } from '$lib/stores/backend';
import { get } from 'svelte/store';
// import { getAccessToken } from '$lib/stores/authentication'; // Not used
import type { PageLoad } from './$types';
// import { authFetch } from '$lib/utilities/authFetch'; // Not used directly here
import { getAllUsers } from '$lib/services/adminDashboard';

// const apiUrl = get(backendUrl); // apiUrl is not used here

export const load = (async ({ fetch }) => {
	// fetch is not used here
	try {
		console.log('Attempting to load users in +page.ts...');
		let { users, error } = await getAllUsers(); // Assuming getAllUsers uses authFetch internally
		if (error) {
			console.error('Error loading users in +page.ts:', error);
			return { users: [], error };
		}
		console.log('Users loaded in +page.ts:', users);
		return { users };
	} catch (error) {
		console.error('Error in +page.ts load function:', error);
		// Depending on how you want to handle errors, you might return an error state
		// or an empty users array to prevent a total crash.
		// For now, re-throwing might show a SvelteKit error page, which is informative.
		// throw error;
		return { users: [], error: 'Failed to load users' }; // Or return an error prop
	}
}) satisfies PageLoad;
