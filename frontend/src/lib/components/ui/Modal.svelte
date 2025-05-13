<script lang="ts">
	import Icon from '@iconify/svelte';

	let {
		isOpen,
		onClose,
		title,
		children
	}: { isOpen: boolean; onClose: () => void; title: string; children: any } = $props();

	console.log('children', typeof children);
</script>

{#if isOpen}
	<div
		class="fixed inset-0 z-50 flex h-screen w-screen items-center justify-center bg-black/50"
		role="dialog"
		aria-modal="true"
		aria-labelledby="modal-title"
	>
		<div class="bg-background-secondary w-full max-w-lg rounded-md p-4 shadow-lg">
			<div class="flex items-center justify-between">
				<h1 id="modal-title" class="text-lg font-bold">{title}</h1>
				<button
					class="text-primary hover:text-primary-focus"
					onclick={onClose}
					aria-label="Close modal"
				>
					<Icon icon="ic:baseline-close" class="text-2xl" />
				</button>
			</div>
			<div class="mt-4">
				{#if children}
					{@render children()}
				{/if}
			</div>
			<div class="mt-4 flex justify-end">
				<button class="button-primary" onclick={onClose}> Close </button>
			</div>
		</div>
	</div>
{/if}
