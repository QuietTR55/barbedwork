<script lang="ts">
    import type { Snippet } from 'svelte';
    import type { LayoutData } from './$types';
	import AdminDashboardSidebar from '$lib/components/ui/AdminDashboardSidebar.svelte';

    let { data, children }: { data: LayoutData, children: Snippet } = $props();
</script>

<div class="bg-background-primary text-text-primary flex flex-col h-screen lg:flex-row md:flex-row"> 
    <div class="w-full lg:w-[198px] md:w-[198px] hidden lg:block md:block">
        <div class="fixed h-screen overflow-y-auto w-[198px]">
            <AdminDashboardSidebar />
        </div>
    </div>
    <div class="bg-background-primary text-text-primary flex flex-col w-full mb-[12px] pb-16 lg:pb-0 md:pb-0 lg:pl-[40px] md:pl-[40px]">
        {@render children()}
    </div>
    <div class="w-full h-12 lg:w-[198px] md:w-[198px] lg:hidden md:hidden fixed bottom-0 left-0">
        <AdminDashboardSidebar />
    </div>
</div>