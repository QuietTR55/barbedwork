<script lang="ts">
	import SingleFileInput from '$lib/components/ui/SingleFileInput.svelte';

	let selectedFile: File | null = null;

	function uploadFile(event: Event) {
		event.preventDefault();
		const input = event.target as HTMLInputElement;
		if (input.files && input.files.length > 0) {
			selectedFile = input.files[0];
			// Handle the file upload logic here
		} else {
			alert('No file selected!');
		}
	}
</script>

<div>
	<form
		class="bg-background-secondary flex w-full flex-col items-center gap-2 rounded-lg p-4 shadow-md md:w-[256px]"
	>
		<SingleFileInput {selectedFile} />
		<input type="text" placeholder="Workspace Name" class="input-primary mb-4 w-full" />
	</form>
</div>
