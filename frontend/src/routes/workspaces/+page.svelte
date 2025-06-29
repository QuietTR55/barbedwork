<script lang="ts">
	import type { Workspace } from '$lib/models/workspace';
	import { get } from 'svelte/store';
	import type { PageProps } from './$types';
	import { backendUrl } from '$lib/stores/backend';
	import Icon from '@iconify/svelte';
	import { goto } from '$app/navigation';
	let { data }: { data: { workspaces: Array<Workspace> } } = $props();
	let url = get(backendUrl);

	function handleWorkspaceClick(workspace: Workspace) {
		goto(`/workspaces/${workspace.Id}`);
	}
</script>

<div class="background-primary-centered">
	<ul class="bg-background-secondary text-text-primary flex w-5/6 flex-col gap-2 rounded-md p-2">
		{#each data.workspaces as workspace}
			<li
				class="bg-background-primary flex flex-row items-center justify-between gap-2 rounded-md p-2"
			>
				<div class="flex min-w-0 flex-row items-center gap-2">
					{#if workspace.ImagePath && workspace.ImagePath.Valid && workspace.ImagePath.String}
						<img
							src="{url}{workspace.ImagePath.String}"
							alt="{workspace.Name} logo"
							class="h-[64px] w-[64px] flex-shrink-0 rounded-full object-cover"
						/>
					{:else}
						<div
							class="bg-background-tertiary flex h-[64px] w-[64px] flex-shrink-0 items-center justify-center rounded-full"
						>
							<Icon icon="ic:baseline-workspaces" class="text-3xl text-gray-400" />
						</div>
					{/if}
					<h2 class="truncate text-lg font-semibold">{workspace.Name}</h2>
				</div>
				<a class="button-primary flex-shrink-0 items-center" href="/workspaces/{workspace.Id}">
					<Icon icon="ic:baseline-log-in" />
					<span class="hidden sm:inline">Open</span>
				</a>
			</li>
		{/each}
	</ul>
</div>
