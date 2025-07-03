<script lang="ts">
	import Icon from '@iconify/svelte';

	let { children } = $props();
	let sidebarOpen = $state(false);

	function toggleSidebar() {
		sidebarOpen = !sidebarOpen;
	}

	function closeSidebar() {
		sidebarOpen = false;
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
		>
			<button
				onclick={closeSidebar}
				class="text-text-primary hover:bg-background-secondary absolute right-4 top-4 rounded-md p-2 transition-colors lg:hidden"
				aria-label="Close sidebar"
			>
				<Icon icon="mdi:close" class="h-5 w-5" />
			</button>

			<div class="mt-12 lg:mt-0">
				<h2 class="text-text-primary mb-4 text-lg font-semibold">Workspace</h2>
				<p class="text-text-secondary">Sidebar content goes here</p>
			</div>
		</div>

		<!-- Main content -->
		<div class="flex-1 p-4 lg:ml-0">
			{@render children()}
		</div>
	</div>
</div>
