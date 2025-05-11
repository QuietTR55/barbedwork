<script lang="ts">
	import Input from '$lib/components/Input.svelte';
	import SingleFileInput from '$lib/components/ui/SingleFileInput.svelte';
	import { backendUrl } from '$lib/stores/backend';
	import { authFetch } from '$lib/utilities/authFetch';
	import { get } from 'svelte/store';
	import type { PageProps } from '../$types';
	import type { Workspace } from '$lib/models/workspace';

	let selectedFile: File | null = null;

	let { data }: { data: { workspaces: Workspace[] } } = $props();
	let workspaces: Workspace[] = $state<Workspace[]>(data.workspaces);
	for (let i = 0; i < workspaces.length; i++) {
		console.log('workspace : ', workspaces[i]);
	}
	let url = get(backendUrl);
	async function createWorkspace(event: Event) {
		event.preventDefault();
		let formData: FormData = new FormData(event.target as HTMLFormElement);
		if (selectedFile) {
			formData.append('image', selectedFile);
		}
		console.log('workspace name : ', formData.get('WorkspaceName'));

		const url = get(backendUrl) + '/api/admin/create-workspace';
		const result = await authFetch(url, {
			method: 'POST',
			body: formData
		});
		if (result.ok) {
			console.log('Workspace created successfully');
			const response = await result.json();

			workspaces = [...workspaces, response.workspace];
		} else {
			console.error('Failed to create workspace:', result.status, await result.text());
		}
	}
</script>

<div class="flex flex-col gap-4 overflow-y-hidden md:flex-row">
	<form
		class="bg-background-secondary flex w-full flex-col items-center gap-2 rounded-lg p-4 shadow-md lg:w-[348px]"
		onsubmit={(event) => {
			createWorkspace(event);
		}}
	>
		<SingleFileInput bind:selectedFile />
		<Input
			icon="ic:sharp-mode-edit"
			type="text"
			placeholder="Workspace Name"
			name="WorkspaceName"
		/>
		<button class="button-primary" type="submit"> Create </button>
	</form>
	{#if workspaces.length > 0}
		<div
			class="bg-background-secondary flex w-full flex-col items-center gap-2 rounded-lg p-4 shadow-md"
		>
			<ul class="flex w-full flex-col gap-2 overflow-auto">
				{#if workspaces.length === 0}
					<li
						class="bg-background-secondary flex flex-row items-center gap-2 rounded-lg p-4 shadow-md"
					>
						No workspaces available
					</li>
				{/if}
				{#if workspaces.length > 0}
					<li
						class="bg-background-secondary flex flex-row items-center gap-2 rounded-lg p-4 shadow-md"
					>
						All workspaces
					</li>
				{/if}
				{#each workspaces as workspace}
					<li
						class="bg-background-secondary flex flex-row items-center gap-2 rounded-lg p-4 shadow-md"
					>
						{#if workspace.ImagePath && workspace.ImagePath.Valid && workspace.ImagePath.String}
							<img
								src={`${url}${workspace.ImagePath.String}`}
								alt={workspace.Name}
								class="h-18 w-18 rounded-full"
								onerror={(event) => {
									if (event.target instanceof HTMLElement) {
										event.target.style.display = 'none';
										if (event.target.nextElementSibling instanceof HTMLElement) {
											event.target.nextElementSibling.style.display = 'flex';
										}
									}
								}}
							/>
							<div
								class="h-18 w-18 flex items-center justify-center rounded-full bg-gray-500 text-xl font-bold text-white"
								style="display: none;"
							>
								{workspace.Name.charAt(0).toUpperCase()}
							</div>
						{:else}
							<div
								class="h-18 w-18 flex items-center justify-center rounded-full bg-gray-500 text-xl font-bold text-white"
							>
								{workspace.Name.charAt(0).toUpperCase()}
							</div>
						{/if}
						<a href={`/admin-dashboard/workspaces/${workspace.Id}`}>
							{workspace.Name}
						</a>
					</li>
				{/each}
			</ul>
		</div>
	{/if}
</div>
