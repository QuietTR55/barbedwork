import type { User } from '$lib/models/user';
import { getAccessToken } from '$lib/stores/authentication';
import { backendUrl } from '$lib/stores/backend';
import { authFetch } from '$lib/utilities/authFetch';
import { get } from 'svelte/store';

const apiUrl = get(backendUrl);

export async function getAllUsers(): Promise<{ users: User[]; error?: Error }> {
	console.log('Fetching users');
	console.log(getAccessToken());
	const response = await authFetch(`${apiUrl}/api/admin/users`, {
		method: 'GET',
		headers: {
			'Content-Type': 'application/json',
			Authorization: `Bearer ${getAccessToken()}`
		},
		credentials: 'include'
	});

	if (!response.ok) {
		console.error('Failed to fetch users:', response.status, await response.text().catch(() => ''));
		return { users: [], error: new Error('Failed to fetch users') };
	}

	const dbUsers: User[] = await response.json();
	console.log('Fetched users:', dbUsers);

	if (Array.isArray(dbUsers)) {
		return { users: dbUsers };
	} else {
		console.error('API response for getAllUsers was not an array as expected:', dbUsers);
		return { users: [], error: new Error('Invalid data format from server for users.') };
	}
}

export async function CreateUser(userName: string, password: string): Promise<User | void> {
	console.log('Creating user:', userName, password);
	const response = await authFetch(`${apiUrl}/api/admin/create-user`, {
		method: 'POST',
		body: JSON.stringify({
			username: userName,
			password: password
		})
	});

	if (!response.ok) {
		console.error('Failed to create user:', response.status, await response.text().catch(() => ''));
		return;
	}

	const apiResponse = await response.json();
	console.log('response', apiResponse);
	if (apiResponse.user) {
		return apiResponse.user;
	} else {
		console.error("API response for create user did not contain a 'user' object:", apiResponse);
		return;
	}
}
