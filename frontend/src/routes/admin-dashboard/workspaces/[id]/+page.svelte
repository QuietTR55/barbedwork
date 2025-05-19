<script lang="ts">
	import type { Workspace } from '$lib/models/workspace';
	import Icon from '@iconify/svelte';
	// import type { PageProps } from '../$types'; // Not used
	import Modal from '$lib/components/ui/Modal.svelte';
	import Input from '$lib/components/Input.svelte';
	// import { each } from 'chart.js/helpers'; // Not used
	import type { User } from '$lib/models/user';
	import { backendUrl } from '$lib/stores/backend';
	import { get } from 'svelte/store';
	import { addUserToWorkspace } from '$lib/services/workspace';

	let { data }: { data: { workspace: Workspace; allUsers: User[] } } = $props();
	console.log('data users', data.workspace.Users);
	let workspace: Workspace = $state<Workspace>(data.workspace);
	let allUsers: User[] = $state<User[]>(data.allUsers);

	console.log('Workspace Users', workspace.Users);

	let addUserModelOn = $state(false);
	let searchTerm = $state('');

	function openModal() {
		searchTerm = '';
		addUserModelOn = true;
	}

	function closeModal() {
		addUserModelOn = false;
	}

	let filteredUsers = $derived(
		allUsers.filter((user) => {
			const isAlreadyInWorkspace = workspace.Users?.some((wsUser) => wsUser.id === user.id);
			if (isAlreadyInWorkspace) {
				return false;
			}
			return user.username.toLowerCase().includes(searchTerm.toLowerCase());
		})
	);

	async function handleAddUser(userId: string) {
		const result = await addUserToWorkspace(workspace.Id, userId);
		if (result) {
			const userToAdd = allUsers.find((user) => user.id === userId);
			if (userToAdd) {
				workspace.Users.push(userToAdd);
			}
		}
	}

	async function handleRemoveUser(userId: string) {
		console.log('removing user with id ', userId);
	}
</script>

<Modal title="Add User to Workspace" isOpen={addUserModelOn} onClose={closeModal}>
	<div class="flex flex-col gap-4">
		<Input
			icon="ic:baseline-search"
			placeholder="Search for user by username"
			bind:value={searchTerm}
			name="userSearch"
		/>
		{#if searchTerm && filteredUsers.length === 0}
			<p class="p-2 text-center text-gray-500">No users found matching "{searchTerm}".</p>
		{:else if !searchTerm && filteredUsers.length === 0 && allUsers.filter((u) => !workspace.Users?.some((wsU) => wsU.id === u.id)).length === 0}
			<p class="p-2 text-center text-gray-500">All users are already in this workspace.</p>
		{:else if filteredUsers.length === 0 && !searchTerm}
			<p class="p-2 text-center text-gray-500">Start typing to search for users to add.</p>
		{/if}

		<ul class="flex flex-col gap-2">
			{#each filteredUsers as user (user.id)}
				<li
					class="bg-background-primary flex w-full items-center justify-between gap-2 rounded-md p-2"
				>
					<div class="flex min-w-0 items-center gap-2">
						{#if user.image_path}
							<!-- TODO implement user profile images -->
						{:else}
							<Icon icon="ic:baseline-person" class="text-primary flex-shrink-0 text-2xl" />
						{/if}
						<p class="truncate">{user.username}</p>
					</div>
					<button class="button-primary flex-shrink-0" onclick={() => handleAddUser(user.id)}>
						<Icon icon="ic:baseline-add" class="text-xl text-white" />
						<span class="hidden sm:inline">Add</span>
					</button>
				</li>
			{/each}
		</ul>
	</div>
</Modal>

<div class="flex flex-col gap-4 overflow-y-auto">
	<div class="bg-background-secondary flex items-center gap-2 rounded-md p-2">
		{#if workspace.ImagePath && workspace.ImagePath.Valid && workspace.ImagePath.String}{:else}
			<Icon icon="ic:baseline-business" class="text-primary text-3xl" />
		{/if}
		<h1 class="truncate text-lg font-bold">{workspace.Name}</h1>
	</div>

	<div class="bg-background-secondary flex flex-col gap-4 overflow-y-auto rounded-md p-2">
		<div class="flex items-center justify-between">
			<h2 class="text-md font-semibold">Users in Workspace</h2>
			<button class="button-primary" onclick={openModal}>
				<Icon icon="ic:baseline-person-add" class="text-xl text-white" />
				<span class="hidden sm:inline">Add User</span>
			</button>
		</div>
		{#if !workspace.Users || workspace.Users.length === 0}
			<div class="flex items-center gap-2 p-2 text-gray-500">
				<Icon icon="ic:baseline-people" class="text-xl" />
				<p>No users currently in this workspace.</p>
			</div>
		{:else}
			<ul class="flex flex-col gap-2">
				{#each workspace.Users as user (user.id)}
					<li class="bg-background-primary flex items-center justify-between gap-2 rounded-md p-2">
						<div class="flex min-w-0 items-center gap-2">
							{#if user.image_path}
								<img
									src={`${get(backendUrl)}${user.image_path}`}
									alt={user.username}
									class="h-8 w-8 flex-shrink-0 rounded-full"
								/>
							{:else}
								<Icon icon="ic:baseline-person" class="text-primary flex-shrink-0 text-xl" />
							{/if}
							<p class="truncate">{user.username}</p>
						</div>
						<button class="button-primary" onclick={() => handleRemoveUser(user.id)}>
							<Icon icon="ic:baseline-remove" class="text-xl text-white" />
							<span class="hidden sm:inline">Remove</span>
						</button>
					</li>
				{/each}
			</ul>
		{/if}
	</div>
	<!-- Other workspace details can go here -->
</div>
