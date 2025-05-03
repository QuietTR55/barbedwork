<script lang="ts">
	import { onMount } from 'svelte';
	import type { PageData } from './$types';
	import { getAccessToken } from '$lib/stores/authentication';
	import { get } from 'svelte/store';
	import { backendUrl } from '$lib/stores/backend';
	import { authFetch } from '$lib/utilities/authFetch';
	import Input from '$lib/components/Input.svelte';
	import Icon from '@iconify/svelte';

	let { data }: { data: PageData } = $props();

	const apiUrl = get(backendUrl);
	let users = $state([]);
	let search = $state('');
	onMount(async () => {
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
			console.error(
				'Failed to fetch users:',
				response.status,
				await response.text().catch(() => '')
			);
			return { status: response.status, error: new Error('Failed to fetch users'), users: [] };
		}

		const dbUsers = await response.json();
		console.log(dbUsers);
		users = dbUsers.users;
	});

	let userName = $state('');
	let password = $state('');

	const createUser = async (event: Event) => {
		event.preventDefault();
		console.log(userName, password);
		const response = await authFetch(`${apiUrl}/api/admin/create-user`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				Authorization: `Bearer ${getAccessToken()}`
			},
			body: JSON.stringify({
				username: userName,
				password: password
			})
		});

		if (!response.ok) {
			console.error(
				'Failed to create user:',
				response.status,
				await response.text().catch(() => '')
			);
			return;
		}

		const data = await response.json();
		console.log(data);
		users.push(data.user);
	};

	// Add this function to truncate text
	function truncateText(text, maxLength = 25) {
		if (text.length <= maxLength) return text;
		return text.substring(0, maxLength) + '...';
	}

	const deactivateUser = async (userid: string) => {
		console.log(userid);
	};
</script>

<div class="flex h-full flex-col gap-4 overflow-y-auto xl:h-max xl:flex-row">
	<form
		class="bg-background-secondary flex w-full flex-col gap-2 rounded-md p-2 xl:h-min xl:w-[400px]"
		onsubmit={createUser}
	>
		<Input
			icon="ic:baseline-person"
			placeholder="UserName"
			type="text"
			bind:value={userName}
			maxlength={50}
		/>
		<Input
			icon="ic:baseline-lock"
			placeholder="Password"
			type="password"
			bind:value={password}
			maxlength={50}
		/>
		<button type="submit" class="button-primary">
			<p class="text-center">Create New User</p>
		</button>
	</form>
	<div class="flex h-full flex-col gap-2 xl:w-full">
		<h1 class="text-2xl font-bold">Users</h1>
		<Input
			icon="ic:baseline-search"
			placeholder="Search"
			type="text"
			bind:value={search}
			maxlength={50}
		/>
		<p class="text-text-secondary text-sm">Showing {users.length} users</p>
		<ul class="flex flex-col gap-2 overflow-y-auto">
			{#each users as user}
				{#if user.username.toLowerCase().includes(search.toLowerCase())}
					<li
						class="bg-background-secondary flex flex-row items-center justify-between rounded-md p-2"
					>
						<p class="truncate">{truncateText(user.username, 50)}</p>
						<button class="button-important" onclick={() => deactivateUser(user.id)}>
							<Icon icon="ic:twotone-no-accounts" />
							<p class="hidden text-center sm:block">Deactivate</p>
						</button>
					</li>
				{/if}
			{/each}
		</ul>
	</div>
</div>
