<script lang="ts">
	import { themeName } from '$lib/stores/theme';
	import { onMount } from 'svelte';
	import '../app.css';
	import ThemeDebugger from '$lib/components/themeDebugger.svelte';

	onMount(() => {
		console.log('onMount');
	});
	
	let currentTheme = $state(themeName);

	let { children } = $props();
</script>

{@render children()}