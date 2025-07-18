<script lang="ts">
	import Input from '$lib/components/Input.svelte';
	import Icon from '@iconify/svelte';
	import type { PageData } from './$types';
	import { onMount } from 'svelte';
	import { backendUrl } from '$lib/stores/backend';
	import { get } from 'svelte/store';
	import { authFetch } from '$lib/utilities/authFetch';
	import { goto } from '$app/navigation';
	import { setAccessToken } from '$lib/stores/authentication';

	let { data }: { data: PageData } = $props();

	let userName = $state('');
	let password = $state('');

	onMount(() => {
		const url = get(backendUrl);
		console.log('Backend URL: ', url);
	});

	const loginUser = async () => {
		const endpoint = `${get(backendUrl)}/api/auth/user-login`;
		const result = await fetch(endpoint, {
			method: 'POST',
			body: JSON.stringify({
				username: userName,
				password: password
			}),
			headers: {
				'Content-Type': 'application/json'
			},

			credentials: 'include'
		});
		if (!result.ok) {
			const errorText = await result.text();
			console.error('Login failed:', result.status, errorText);
			return;
		}
		let accessToken: string | null = null;
		try {
			const data = await result.json();
			accessToken = data.accessToken;
		} catch (err) {
			console.error('Error parsing JSON:', err);
			return;
		}
		if (!accessToken) {
			console.error('No access token received');
			return;
		}
		console.log('Access token received:', accessToken);
		setAccessToken(accessToken);

		if (result.ok) {
			goto('/workspaces');
		}
	};
</script>

<div class="background-primary-centered flex flex-col items-center justify-center gap-2">
	<Icon icon="material-symbols:person" class="text-accent mb-4 w-full text-left text-6xl" />
	<h1 class="text-text-primary w-full text-center text-2xl font-bold">User Login</h1>
	<form
		class="bg-background-secondary text-text-primary mt-4 flex w-[90%] flex-col items-center justify-center gap-4 rounded-md p-4 sm:w-6/12 md:w-4/12 lg:w-3/12"
	>
		<p class="text-text-secondary w-full text-left">Username</p>
		<Input
			icon="material-symbols:person"
			placeholder="Username"
			type="text"
			bind:value={userName}
		/>
		<p class="text-text-secondary w-full text-left">Password</p>
		<Input
			icon="material-symbols:lock"
			placeholder="Password"
			type="password"
			bind:value={password}
		/>
		<button class="button-primary" onclick={loginUser}>Login</button>
	</form>
	<p class="w-full p-4 text-center">
		Contact your workspace administrator for an invitation to join.
	</p>
	<p class="text-text-secondary mt-4 text-center">
		Forgot your password?
		<a href="/reset-password" class="text-accent">Reset it here</a>.
	</p>
</div>
