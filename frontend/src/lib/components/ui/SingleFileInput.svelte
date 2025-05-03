<script lang="ts">
	import Icon from '@iconify/svelte';
	let { selectedFile = $bindable(null) }: { selectedFile: File | null } = $props();
	const fileInput = document.createElement('input');
	fileInput.type = 'file';
	fileInput.multiple = false;
	fileInput.accept = 'image/*';
	fileInput.onchange = handleFileChange;
	const state = $state({
		files: [] as File[]
	});

	function handleFileChange(event: Event) {
		const input = event.target as HTMLInputElement;
		const selected = Array.from(input.files ?? []);
		state.files = selected;
		selectedFile = selected[0] || null; // Update the bound file
	}

	function removeImage(index: number) {
		state.files.splice(index, 1);
		state.files = [...state.files]; // ensure reactivity
	}

	function getImageURL(file: File) {
		return URL.createObjectURL(file);
	}

	function triggerFileInput() {
		fileInput.click();
	}
</script>

<div class="flex flex-col gap-4">
	<!-- Clickable box to trigger image input -->
	<div
		onclick={() => {
			triggerFileInput();
		}}
		onkeydown={() => {
			triggerFileInput();
		}}
		class="flex h-32 w-32 cursor-pointer items-center justify-center rounded border-2 border-dashed border-gray-300 transition-colors hover:border-blue-500"
	>
		{#if state.files.length === 0}
			<Icon icon="ic:baseline-image" class="text-4xl text-gray-400" />
		{:else}
			<div class="relative flex h-32 w-32 overflow-hidden rounded">
				<img
					src={getImageURL(state.files[0])}
					alt="Selected"
					class="h-full w-full rounded object-cover object-center"
				/>
				<button
					onclick={(event) => {
						event.stopPropagation();
						removeImage(0);
					}}
					class="bg-accent/60 hover:bg-accent/80 absolute top-1 right-1 flex h-6 w-6 items-center justify-center rounded-full text-white"
					aria-label="Remove image"
				>
					x
				</button>
			</div>
		{/if}
	</div>
</div>
