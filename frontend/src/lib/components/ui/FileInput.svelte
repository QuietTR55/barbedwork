<script lang="ts">
	import { onDestroy } from 'svelte';

	let {
		icon,
		file = $bindable(null),
		multiple = false,
		accept = 'image/*',
		previewSize = '100px',
		...rest
	} = $props();

	let fileInput: HTMLInputElement;
	let previews: string[] = $state([]);

	function handleChange(event: Event) {
		const input = event.target as HTMLInputElement;
		if (!input.files) {
			return;
		}

		// Clear previous object URLs to prevent memory leaks
		clearPreviews();

		file = multiple ? input.files : input.files[0];
		createPreviews();
	}

	function createPreviews() {
		if (!file) return;

		if (multiple && file instanceof FileList) {
			previews = Array.from(file).map((f) => URL.createObjectURL(f));
		} else if (file instanceof File) {
			previews = [URL.createObjectURL(file)];
		}
	}

	function clearPreviews() {
		previews.forEach((url) => URL.revokeObjectURL(url));
		previews = [];
	}

	function removePreview(index: number) {
		URL.revokeObjectURL(previews[index]);
		previews.splice(index, 1);
		previews = [...previews]; // Trigger reactivity
	}

	onDestroy(clearPreviews);
</script>

<div class="flex flex-col gap-3">
	<input
		type="file"
		bind:this={fileInput}
		{multiple}
		{accept}
		onchange={handleChange}
		class="sr-only"
		{...rest}
	/>

	<button
		type="button"
		onclick={() => fileInput.click()}
		class="bg-background-secondary border-border-primary hover:bg-accent flex items-center justify-center gap-2 rounded border px-4 py-2 text-sm"
	>
		{#if icon}<span>{icon}</span>{/if}
		<span>Select {multiple ? 'Images' : 'Image'}</span>
	</button>

	{#if previews.length > 0}
		<div class="flex flex-wrap gap-2">
			{#each previews as preview, index}
				<div
					class="relative overflow-hidden rounded border border-slate-200"
					style="width: {previewSize}; height: {previewSize};"
				>
					<button
						type="button"
						class="absolute top-1 right-1 flex h-5 w-5 items-center justify-center rounded-full bg-red-500 text-xs text-white"
						onclick={() => removePreview(index)}
					>
						&times;
					</button>
					<img src={preview} alt="Preview" class="h-full w-full object-cover" />
				</div>
			{/each}
		</div>
	{/if}
</div>
