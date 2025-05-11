<script lang="ts">
	import type { Workspace } from '$lib/models/workspace';
	import type { PageProps } from '../$types';

	let { data }: { data: { workspace: Workspace } } = $props();
	let workspace: Workspace = $state<Workspace>(data.workspace);
</script>

<div>
	<h1>{workspace.Name}</h1>
	<ul></ul>
</div>
