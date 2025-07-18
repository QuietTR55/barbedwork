<script lang="ts">
	import type { Workspace } from '$lib/models/workspace';
	import { createChannel } from '$lib/services/workspace';
	import { activeWorkspace, workspaceActions } from '$lib/stores/workspace';
	import Icon from '@iconify/svelte';
	import { get } from 'svelte/store';

	let { data, children }: { data: { workspace: Workspace }; children: any } = $props();
	let sidebarOpen = $state(false);
	let contextMenuOpen = $state(false);
	let contextMenuX = $state(0);
	let contextMenuY = $state(0);
	let showChannelForm = $state(false);
	let channelName = $state('');
	let selectedEmoji = $state('💬');

	const commonEmojis = [
		'💬',
		'📢',
		'🔊',
		'💡',
		'🎯',
		'🚀',
		'📝',
		'🎮',
		'🎵',
		'📚',
		'⚡',
		'🔥',
		'🌟',
		'❤️',
		'🎉'
	];

	function toggleSidebar() {
		sidebarOpen = !sidebarOpen;
	}

	function closeSidebar() {
		sidebarOpen = false;
	}

	function handleRightClick(event: MouseEvent) {
		event.preventDefault();
		contextMenuX = event.clientX;
		contextMenuY = event.clientY;
		contextMenuOpen = true;
	}

	function closeContextMenu() {
		contextMenuOpen = false;
	}

	function openChannelForm() {
		showChannelForm = true;
		contextMenuOpen = false;
		channelName = '';
		selectedEmoji = '💬';
	}

	function closeChannelForm() {
		showChannelForm = false;
		channelName = '';
		selectedEmoji = '💬';
	}

	async function createNewChannel() {
		if (!channelName.trim()) {
			console.error('Channel name is required');
			return;
		}

		console.log('Adding new channel...', { name: channelName, emoji: selectedEmoji });
		const createdChannel = await createChannel(get(activeWorkspace)?.Id ?? '', {
			name: channelName.trim(),
			emoji: selectedEmoji
		});

		if (!createdChannel) {
			console.error('Failed to create channel');
			return;
		}

		workspaceActions.addChannel(createdChannel);
		closeChannelForm();
	}
</script>

<div class="background-primary-centered p-2">
	<div class="bg-background-secondary relative flex h-full w-full rounded-md">
		<button
			onclick={toggleSidebar}
			class="bg-background-tertiary text-text-primary hover:bg-background-secondary fixed left-4 top-4 z-50 rounded-md p-2 transition-colors lg:hidden"
			aria-label="Toggle sidebar"
		>
			<Icon icon="mdi:menu" class="h-6 w-6" />
		</button>

		{#if sidebarOpen}
			<div
				class="fixed inset-0 z-30 bg-black bg-opacity-50 lg:hidden"
				onclick={closeSidebar}
				onkeydown={(e) => e.key === 'Escape' && closeSidebar()}
				role="button"
				tabindex="0"
				aria-label="Close sidebar"
			></div>
		{/if}

		<div
			class="
			{sidebarOpen ? 'translate-x-0' : '-translate-x-full'}
			bg-background-tertiary
			fixed left-0
			top-0 z-40
			h-full w-80
			rounded-l-md
			p-2
			transition-transform
			duration-300 ease-in-out lg:static
			lg:z-auto
			lg:translate-x-0
		"
			oncontextmenu={handleRightClick}
			role="navigation"
			aria-label="Workspace sidebar"
		>
			<button
				onclick={closeSidebar}
				class="text-text-primary hover:bg-background-secondary absolute right-4 top-4 rounded-md p-2 transition-colors lg:hidden"
				aria-label="Close sidebar"
			>
				<Icon icon="mdi:close" class="h-5 w-5" />
			</button>

			<div class="mt-12 lg:mt-0">
				<h2 class="text-text-primary mb-4 text-lg font-semibold">{data.workspace.Name}</h2>

				<!-- Channels section -->
				<div class="min-h-0 flex-1">
					<h3 class="text-text-primary mb-2 text-sm font-medium uppercase tracking-wide">
						Channels
					</h3>
					<div class="h-full">
						{#each data.workspace.Channels as channel}
							<div
								class="text-text-secondary hover:text-text-primary hover:bg-background-secondary mb-1 flex cursor-pointer items-center gap-2 rounded px-2 py-1 transition-colors"
							>
								<span>{channel.emoji}</span>
								<span>{channel.name}</span>
							</div>
						{/each}
					</div>
				</div>
			</div>
		</div>

		<!-- Main content -->
		<div class="flex-1 p-4 lg:ml-0">
			{@render children()}
		</div>

		<!-- Context Menu -->
		{#if contextMenuOpen}
			<div
				class="bg-background-tertiary border-border-primary fixed z-50 min-w-48 rounded-md border shadow-lg"
				style="left: {contextMenuX}px; top: {contextMenuY}px;"
			>
				<div class="p-1">
					<button
						onclick={openChannelForm}
						class="text-text-primary hover:bg-background-secondary flex w-full items-center gap-2 rounded px-3 py-2 text-left text-sm transition-colors"
					>
						<Icon icon="mdi:plus" class="h-4 w-4" />
						Add Channel
					</button>
				</div>
			</div>
		{/if}
	</div>

	<!-- Channel Creation Form Modal -->
	{#if showChannelForm}
		<div class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
			<div
				class="bg-background-secondary border-border-primary mx-4 w-full max-w-md rounded-lg border p-6 shadow-xl"
			>
				<h3 class="text-text-primary mb-4 text-lg font-semibold">Create New Channel</h3>

				<form
					onsubmit={(e) => {
						e.preventDefault();
						createNewChannel();
					}}
					class="space-y-4"
				>
					<!-- Channel Name Input -->
					<div>
						<label for="channel-name" class="text-text-primary mb-2 block text-sm font-medium">
							Channel Name
						</label>
						<input
							id="channel-name"
							type="text"
							bind:value={channelName}
							placeholder="Enter channel name"
							class="bg-background-tertiary border-border-primary text-text-primary placeholder:text-text-tertiary w-full rounded-md border px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
							required
						/>
					</div>

					<!-- Emoji Selection -->
					<div>
						<div class="text-text-primary mb-2 block text-sm font-medium">Channel Emoji</div>
						<div class="grid grid-cols-5 gap-2">
							{#each commonEmojis as emoji}
								<button
									type="button"
									onclick={() => (selectedEmoji = emoji)}
									class="hover:bg-background-tertiary border-border-primary flex h-10 w-10 items-center justify-center rounded-md border text-lg transition-colors {selectedEmoji ===
									emoji
										? 'border-blue-500 bg-blue-500'
										: ''}"
								>
									{emoji}
								</button>
							{/each}
						</div>
						<div class="mt-2">
							<input
								type="text"
								bind:value={selectedEmoji}
								placeholder="Or enter custom emoji"
								class="bg-background-tertiary border-border-primary text-text-primary placeholder:text-text-tertiary w-full rounded-md border px-3 py-2 text-center focus:outline-none focus:ring-2 focus:ring-blue-500"
								maxlength="2"
							/>
						</div>
					</div>

					<!-- Form Actions -->
					<div class="flex justify-end gap-3 pt-4">
						<button
							type="button"
							onclick={closeChannelForm}
							class="bg-background-tertiary text-text-secondary hover:bg-background-primary rounded-md px-4 py-2 transition-colors"
						>
							Cancel
						</button>
						<button
							type="submit"
							class="rounded-md bg-blue-600 px-4 py-2 text-white transition-colors hover:bg-blue-700"
						>
							Create Channel
						</button>
					</div>
				</form>
			</div>
		</div>
	{/if}

	<!-- Click outside handler for context menu -->
	{#if contextMenuOpen}
		<div
			class="fixed inset-0 z-40"
			role="button"
			tabindex="0"
			aria-label="Close context menu"
			onclick={closeContextMenu}
			oncontextmenu={(e) => {
				e.preventDefault();
				closeContextMenu();
			}}
			onkeydown={(e) => {
				if (e.key === 'Enter' || e.key === ' ') {
					e.preventDefault();
					closeContextMenu();
				}
			}}
		></div>
	{/if}
</div>
