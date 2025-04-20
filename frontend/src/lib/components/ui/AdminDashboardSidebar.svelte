<script lang="ts">
    import Icon from '@iconify/svelte';
    import { goto } from '$app/navigation';
	import Button from './Button.svelte';

    const widthStyling = "w-full lg:w-[198px] md:w-[198px]";
</script>

<div class="admin-dashboard-side-panel">
    <h1 class="text-2xl font-bold text-center hidden lg:block md:block">Admin Dashboard</h1>
    <Button icon="ri:dashboard-fill" onclick={() => goto('/admin-dashboard')}>
        <p>Home</p>
    </Button>
    <Button icon="material-symbols:person" onclick={() => goto('/admin-dashboard/users')}>
        <p>Users</p>
    </Button>
</div>