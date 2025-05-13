<script lang="ts">
	import { onMount } from 'svelte';
	import type { PageData } from './$types';
	import { getAccessToken } from '$lib/stores/authentication';
	import { get } from 'svelte/store';
	import { backendUrl } from '$lib/stores/backend';
	import { authFetch } from '$lib/utilities/authFetch';
	import Input from '$lib/components/Input.svelte';
	import Icon from '@iconify/svelte';
	import type { User } from '$lib/models/user';
	import { CreateUser } from '$lib/services/adminDashboard';

	const apiUrl = get(backendUrl);
	let { data: pageData }: { data: { users: User[] } } = $props();
	let search = $state('');

	let userName = $state('');
	let password = $state('');

	const createUser = async (event: Event) => {
		event.preventDefault();
		console.log(userName, password);
		let createdUser = await CreateUser(userName, password);
		if (createdUser) {
			pageData.users = [...pageData.users, createdUser];
			userName = '';
			password = '';
		} else {
			console.error('Failed to create user');
		}
	};

	function truncateText(text: string, maxLength = 25) {
		if (text.length <= maxLength) return text;
		return text.substring(0, maxLength) + '...';
	}

	const deactivateUser = async (userid: string) => {
		console.log(userid);
	};
</script>

<div class="flex h-full flex-col gap-4 overflow-y-auto lg:h-max lg:flex-row">
	<form
		class="bg-background-secondary flex w-full flex-col gap-2 rounded-md p-2 lg:h-min lg:w-[400px]"
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
	<div class="flex h-full flex-col gap-2 lg:w-full">
		<h1 class="text-2xl font-bold">Users</h1>
		<Input
			icon="ic:baseline-search"
			placeholder="Search"
			type="text"
			bind:value={search}
			maxlength={50}
		/>
		<p class="text-text-secondary text-sm">Showing {pageData.users.length} users</p>
		<ul class="flex flex-col gap-2 overflow-y-auto">
			{#each pageData.users as user (user.id)}
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
