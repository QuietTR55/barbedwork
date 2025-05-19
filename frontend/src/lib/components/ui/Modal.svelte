<script lang="ts">
	import Icon from '@iconify/svelte';

	let {
		isOpen,
		onClose,
		title,
		children // children prop is used by {@render children()}
	}: { isOpen: boolean; onClose: () => void; title: string; children: any } = $props();

	// $inspect(isOpen); // Optional: for debugging if isOpen changes
</script>

{#if isOpen}
	<div
		class="fixed inset-0 z-50 flex h-screen w-screen items-center justify-center bg-black/50 p-4"
		role="dialog"
		aria-modal="true"
		aria-labelledby="modal-title"
		tabindex="0"
		onclick={(event) => {
			// Optional: Close modal if backdrop is clicked
			if (event.target === event.currentTarget) {
				onClose();
			}
		}}
		onkeydown={(event) => {
			if (event.key === 'Escape' && event.target === event.currentTarget) {
				onClose();
			}
		}}
	>
		<div
			class="bg-background-secondary flex max-h-[90vh] w-full max-w-lg flex-col rounded-md shadow-lg"
		>
			<!-- Header -->
			<div class="border-border-primary flex-shrink-0 border-b p-4">
				<div class="flex items-center justify-between">
					<h1 id="modal-title" class="text-lg font-bold">{title}</h1>
					<button
						class="text-primary hover:text-primary-focus -m-2 p-2"
						onclick={onClose}
						aria-label="Close modal"
					>
						<Icon icon="ic:baseline-close" class="text-2xl" />
					</button>
				</div>
			</div>

			<!-- Content Body - This part will scroll -->
			<div class="min-h-0 flex-grow overflow-y-auto p-4">
				{#if children}
					{@render children()}
				{/if}
			</div>

			<!-- Footer -->
			<div class="flex-shrink-0 border-t border-gray-700 p-4">
				<div class="flex justify-end">
					<button class="button-primary" onclick={onClose}> Close </button>
				</div>
			</div>
		</div>
	</div>
{/if}
