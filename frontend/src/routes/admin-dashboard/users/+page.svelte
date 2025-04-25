<script lang="ts">
	import { onMount } from 'svelte';
    import type { PageData } from './$types';
	import { getAccessToken } from '$lib/stores/authentication';
	import { get } from 'svelte/store';
	import { backendUrl } from '$lib/stores/backend';
	import { authFetch } from '$lib/utilities/authFetch';
	import Input from '$lib/components/Input.svelte';

    let { data }: { data: PageData } = $props();

    const apiUrl = get(backendUrl);

    onMount(async () => {
        console.log('Fetching users');
	    console.log(getAccessToken());
	    const response = await authFetch(`${apiUrl}/admin/users`, {
            method: 'GET',
	    	headers: {
	    		'Content-Type': 'application/json',
	    		'Authorization': `Bearer ${getAccessToken()}`
	    	},
	    	credentials: 'include'
	    });

	    if (!response.ok) {
	    	console.error('Failed to fetch users:', response.status, await response.text().catch(() => ''));
	    	return { status: response.status, error: new Error('Failed to fetch users'), users: [] };
	    }

		const users = await response.json();
		console.log(users);
	});

	let userName = $state('');
	let password = $state('');

	const createUser = async (event: Event) => {
		event.preventDefault();
		console.log(userName, password);
		const response = await authFetch(`${apiUrl}/admin/create-user`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				'Authorization': `Bearer ${getAccessToken()}`
			},
			body: JSON.stringify({
				username: userName,
				password: password
			})
		});

		if (!response.ok) {
			console.error('Failed to create user:', response.status, await response.text().catch(() => ''));
			return;
		}
		
		const data = await response.json();
		console.log(data);
	}
</script>

<div class="flex flex-col gap-4 p-4">
	<form class="w-full bg-background-secondary rounded-md p-2 gap-2 flex flex-col" onsubmit={createUser}>
		<Input icon="ic:baseline-person" placeholder="UserName" type="text" bind:value={userName}/>
		<Input icon="ic:baseline-lock" placeholder="Password" type="password" bind:value={password}/>
		<button type="submit" class="button-primary">
			<p class="text-center">Create New User</p>
		</button>
	</form>
    
</div>