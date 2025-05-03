<script>
	import Icon from '@iconify/svelte';
	import Input from '$lib/components/Input.svelte';
	import { backendUrl } from '$lib/stores/backend';
	import { get } from 'svelte/store';
	import { setAccessToken } from '$lib/stores/authentication';
	import { goto } from '$app/navigation';
	import { onMount } from 'svelte';

	let secretKey = $state('');

	onMount(() => {
		const url = get(backendUrl);
		console.log('Backend URL: ', url);
	});

	async function handleSubmit() {
		try {
			const url = get(backendUrl);
			const response = await fetch(url + '/api/auth/admin/login', {
				// Updated endpoint to include /login
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ secretKey }),
				credentials: 'include'
			});
			console.log('Response', response); // Log the response status
			if (response.ok) {
				// Handle successful login
				console.log('Login successful:', response);
				// Added response for better logging
				const data = await response.json();
				if (data) {
					if (data.accessToken) {
						// Check if accessToken is present
						console.log('Access token received:', data.accessToken); // Log the access token
						setAccessToken(data.accessToken);
						goto('/admin-dashboard'); // Redirect to admin dashboard
					} else {
						console.error('Access token not found in response');
					}
				}
			} else {
				// Handle login error
				console.error('Login failed');
				return response.json().then((errorData) => {
					console.error('Error details:', errorData); // Log error details
				});
			}
		} catch (error) {
			console.error('Error during login:', error); // Log any errors that occur during the fetch
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
