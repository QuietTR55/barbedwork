<script lang="ts">
	import type { Workspace } from '$lib/models/workspace';
	import Icon from '@iconify/svelte';
	import type { PageProps } from '../$types';
	import Modal from '$lib/components/ui/Modal.svelte';
	import Input from '$lib/components/Input.svelte';

	let { data }: { data: { workspace: Workspace } } = $props();
	let workspace: Workspace = $state<Workspace>(data.workspace);

	let addUserModelOn = $state(false);

	console.log('users: ', workspace.users);
	function openModal() {
		console.log('open modal');
		addUserModelOn = true;
	}

	function closeModal() {
		addUserModelOn = false;
	}
</script>

<Modal title="Add User" isOpen={addUserModelOn} onClose={closeModal}>
	<div class="flex flex-col gap-4">
		<Input icon="ic:baseline-email" placeholder="Enter email" />
	</div>
</Modal>
<div class="flex flex-col gap-4 overflow-y-auto">
	<div class="bg-background-secondary flex items-center gap-2 rounded-md p-2">
		<Icon icon="ic:baseline-business" class="text-primary text-2xl" />
		<h1 class="truncate text-lg font-bold">{workspace.Name}</h1>
	</div>
	<div class="bg-background-secondary flex flex-col gap-4 overflow-y-auto rounded-md p-2">
		<div class="flex items-center justify-between">
			<p>Users</p>
			<button class="button-primary" onclick={openModal}>
				<Icon icon="ic:baseline-person-add" class="text-2xl text-white" />
				Add User
			</button>
		</div>
		<ul>
			{#if !workspace.users}
				<li class="flex items-center gap-2">
					<Icon icon="ic:baseline-person-off" class="text-primary text-2xl" />
					<p class="truncate">No users found</p>
				</li>
			{/if}
			{#each workspace.users as user}
				<li class="flex items-center gap-2">
					<Icon icon="ic:baseline-person" class="text-primary text-2xl" />
					<p class="truncate">{user.username}</p>
				</li>
			{/each}
		</ul>
	</div>
</div>
