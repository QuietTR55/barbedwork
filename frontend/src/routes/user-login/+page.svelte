<script lang="ts">
	import Input from '$lib/components/Input.svelte';
	import Icon from '@iconify/svelte';
    import type { PageData } from './$types';

    let { data }: { data: PageData } = $props();
</script>

<div class="background-primary-centered flex flex-col items-center justify-center gap-4">
    <Icon icon="material-symbols:person" class="text-accent text-left w-full text-6xl mb-4"/>
    <h1 class="text-text-primary text-center font-bold w-full text-2xl mb-4">User Login</h1>
    <p class="text-text-secondary text-center">Enter your credentials to login</p>
    <form class="bg-background-secondary w-1/3 flex flex-col items-center justify-center gap-4 rounded-md p-4 text-text-primary mt-4">
        <p class="text-text-secondary text-left w-full">Username</p>
        <Input icon="material-symbols:person" placeholder="Username" type="text"/>
        <p class="text-text-secondary text-left w-full">Password</p>
        <Input icon="material-symbols:lock" placeholder="Password" type="password"/>
        <button class="button-primary">Login</button>
    </form>
    <p>Contact your workspace administrator for an invitation to join.</p>
    <p class="text-text-secondary text-center mt-4">Forgot your password? 
        <a href="/reset-password" class="text-accent">Reset it here</a>.
    </p>
</div>