<script lang="ts">
	import Icon from '@iconify/svelte';
	let { icon, value = $bindable(''), ...props } = $props();
</script>

<div class="relative w-full">
	<div
		class="absolute top-0 bottom-0 left-0 flex items-center justify-center pl-3"
		style="height: 100%;"
	>
		<Icon {icon} class="text-textSecondary h-5 w-5" />
	</div>
	<input class="input-primary pl-10" bind:value {...props} />
</div>

