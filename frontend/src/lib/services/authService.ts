import { goto } from "$app/navigation";
import { setAccessToken } from "$lib/stores/authentication";
import { backendUrl } from "$lib/stores/backend";
import { get } from "svelte/store";

export async function adminLogin(secretKey: string): Promise<boolean> {
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
					} else {
						console.error('Access token not found in response');
					}
				}
			} else {
				// Handle login error
				console.error('Login failed');
			}
            return response.ok;
}