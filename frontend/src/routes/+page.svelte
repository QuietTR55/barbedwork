<script>
	import { goto } from "$app/navigation";
	import Input from "$lib/components/Input.svelte";
	import BarbedWorkLogo from "$lib/components/ui/BarbedWorkLogo.svelte";
	import { backendUrl, setBackendUrl } from "$lib/stores/backend";
	import { get } from "svelte/store";

    let inputUrl = $state(get(backendUrl))
    let error = $state("")

    function updateUrl() {
        setBackendUrl(inputUrl)
    }

    function checkBackend() {
        error = ""
        fetch(inputUrl + "/api/health", {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
            },
        })
            .then((response) => {
                if (response.ok) {
                    console.log("Backend is reachable");
                    setBackendUrl(inputUrl)
                    goto("/user-login")
                } else {
                    console.error("Backend is not reachable");
                    error = "Backend is not reachable";
                }
            })
            .catch((error) => {
                console.error("Error fetching backend:", error);
                error = "Error fetching backend: " + error;
            });
    }

</script>

<div class="background-primary-centered flex flex-col gap-2">
    <p class="flex flex-row items-center gap-2 text-2xl font-bold">
        <BarbedWorkLogo color="#8b5cf6"/>
        <span class=" bg-gradient-to-r from-violet-500 to-violet-400 bg-clip-text text-transparent">BarbedWork</span>
    </p>
    <div class="bg-background-secondary rounded-md w-[90%] sm:w-[600px] text-centered p-2">
        <p>Enter your workspace url</p>
        <Input placeholder="Workspace URL" icon="humbleicons:url" bind:value={inputUrl}/>
    </div>
    <button class="button-primary" onclick={() => {updateUrl() ; checkBackend()}}>
        <span class="text-sm font-medium">Connect</span>
    </button>
    {#if error}
        <p class="text-red-500 text-sm">{error}</p>
    {/if}
</div>