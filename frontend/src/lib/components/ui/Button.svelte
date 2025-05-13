<script lang="ts">
	import Icon from '@iconify/svelte';
	import type { Snippet } from 'svelte';

	let { children, icon, ...props } = $props();
	console.log('children', typeof children);
</script>

<button class="button-primary inline-flex items-center gap-2" {...props}>
	<Icon {icon} class="h-5 w-5 flex-shrink-0" />
	<span class="hidden sm:inline"> {@render children()} </span>
</button>
