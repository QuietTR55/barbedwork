<script>
    import Icon from "@iconify/svelte";
    import Input from "$lib/components/Input.svelte";
	import { backendUrl } from "$lib/stores/backend";
	import { get } from "svelte/store";
	import { setAccessToken } from "$lib/stores/authentication";
	import { goto } from "$app/navigation";
	import { onMount } from "svelte";

    let secretKey = $state("");

    onMount(()=>{
        const url = get(backendUrl);
        console.log("Backend URL: ", url);
    });

    function handleSubmit() {
        const url = get(backendUrl);
        fetch(url + "/api/auth/admin/login", {  // Updated endpoint to include /login
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ secretKey }),
            credentials: "include"
        })
            .then((response) => {
                console.log("Response", response); // Log the response status
                if (response.ok) {
                    // Handle successful login
                    console.log("Login successful:", response);
                    // Added response for better logging
                    response.json().then((json) => {
                        console.log("access token : ", json["accessToken"]);
                        setAccessToken(json["accessToken"]);
                        goto("/admin-dashboard");
                    });
                } else {
                    // Handle login error
                    console.error("Login failed");
                    return response.json().then((errorData) => {
                        console.error("Error details:", errorData); // Log error details
                    });
                }
            })
            .catch((error) => {
                console.error("Error during login:", error);
            });
    }
</script>

<div class="background-primary-centered">
    <Icon icon="eos-icons:admin" class="text-accent text-left w-full text-6xl mb-4"/>
    <h1 class="text-text-primary text-center font-bold w-full text-2xl mb-4">Admin Dashboard Login</h1>
    <p class="text-text-secondary text-center w-full text-sm mb-4">Enter your credentials to login</p>
    <div class="flex flex-col items-center justify-center bg-background-secondary p-4 rounded-md gap-2 w-1/4">
        <p class="text-text-secondary text-left w-full">Secret key</p>
        <Input icon="material-symbols:lock" placeholder="Secret key" bind:value={secretKey} type="password"/>
        <button type="button" class="button-primary" onclick={handleSubmit}>Login</button>
    </div>
</div>