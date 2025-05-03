<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { LayoutData } from './$types';
	import AdminDashboardSidebar from '$lib/components/ui/AdminDashboardSidebar.svelte';

	let { data, children }: { data: LayoutData; children: Snippet } = $props();
</script>

<div class="bg-background-primary text-text-primary flex h-screen flex-col md:flex-row lg:flex-row">
	<div class="hidden md:block md:w-[198px] lg:block lg:w-[198px]">
		<div class="fixed h-screen w-[198px] overflow-y-auto">
			<AdminDashboardSidebar />
		</div>
	</div>
	<div
		class="bg-background-primary text-text-primary mb-[12px] flex w-full flex-col p-2 pb-[50px] md:pb-0 md:pl-[50px]"
	>
		{@render children()}
	</div>
	<div class="fixed bottom-0 left-0 h-[48px] w-full md:hidden md:w-[198px] lg:w-[198px]">
		<AdminDashboardSidebar />
	</div>
</div>

