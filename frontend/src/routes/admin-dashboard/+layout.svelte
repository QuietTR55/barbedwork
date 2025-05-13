<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { LayoutData } from './$types';
	import AdminDashboardSidebar from '$lib/components/ui/AdminDashboardSidebar.svelte';

	let { data, children }: { data: LayoutData; children: Snippet } = $props();
</script>

<div class="bg-background-primary text-text-primary flex h-screen flex-col lg:flex-row">
	<div class="hidden lg:block lg:w-[198px]">
		<div class="fixed h-screen w-[198px] overflow-y-auto">
			<AdminDashboardSidebar />
		</div>
	</div>
	<div
		class="bg-background-primary text-text-primary mb-[12px] flex w-full flex-col p-2 pb-[50px] lg:min-w-0 lg:flex-1 lg:pb-0"
	>
		{@render children()}
	</div>
	<div class="fixed bottom-0 left-0 h-[48px] w-full lg:hidden lg:w-[198px]">
		<AdminDashboardSidebar />
	</div>
</div>
