<script lang="ts">
	import Icon from '@iconify/svelte';
</script>

<div class="admin-dashboard-side-panel">
	<h1 class="hidden text-center text-2xl font-bold md:block lg:block">Admin Dashboard</h1>
	<a href="/admin-dashboard" class="button-primary">
		<Icon icon="ic:baseline-dashboard" />
		<p class="hidden md:block">Dashboard</p>
	</a>
	<a href="/admin-dashboard/users" class="button-primary">
		<Icon icon="ic:baseline-person" />
		<p class="hidden md:block">Users</p>
	</a>
	<a href="/admin-dashboard/workspaces" class="button-primary">
		<Icon icon="ic:twotone-workspaces" />
		<p class="hidden md:block">Workspaces</p>
	</a>
</div>
