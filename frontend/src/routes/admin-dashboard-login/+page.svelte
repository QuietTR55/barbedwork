<script>
	import Icon from '@iconify/svelte';
	import Input from '$lib/components/Input.svelte';
	import { backendUrl } from '$lib/stores/backend';
	import { get } from 'svelte/store';
	import { setAccessToken } from '$lib/stores/authentication';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';
	import { adminLogin } from '$lib/services/authService';

	let secretKey = $state('');

	onMount(() => {
		const url = get(backendUrl);
		console.log('Backend URL: ', url);
	});

	async function handleSubmit() {
		const loggedIn = await adminLogin(secretKey);
		if (loggedIn) {
			goto('/admin-dashboard'); // Redirect to admin dashboard
		}
	}
</script>

<div class="background-primary-centered">
	<Icon icon="eos-icons:admin" class="text-accent mb-4 w-full text-left text-6xl" />
	<h1 class="text-text-primary mb-4 w-full text-center text-2xl font-bold">
		Admin Dashboard Login
	</h1>
	<p class="text-text-secondary mb-4 w-full text-center text-sm">Enter your credentials to login</p>
	<div
		class="bg-background-secondary flex w-1/4 flex-col items-center justify-center gap-2 rounded-md p-4"
	>
		<p class="text-text-secondary w-full text-left">Secret key</p>
		<Input
			icon="material-symbols:lock"
			placeholder="Secret key"
			bind:value={secretKey}
			type="password"
		/>
		<button type="button" class="button-primary" onclick={handleSubmit}>Login</button>
	</div>
</div>
