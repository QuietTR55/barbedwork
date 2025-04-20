<script lang="ts">
	import Input from '$lib/components/Input.svelte';
	import Icon from '@iconify/svelte';
    import type { PageData } from './$types';

    let { data }: { data: PageData } = $props();
</script>

<div class="background-primary-centered">
    <Icon icon="material-symbols:person" class="text-accent text-left w-full text-6xl"/>
    <h1 class="text-2xl font-bold text-text-primary">Login</h1>
    <p class="text-text-secondary text-center">Enter your credentials to login</p>
    <form class="bg-background-secondary w-1/3 flex flex-col items-center justify-center gap-4 rounded-md p-4 text-text-primary mt-4">
        <p class="text-text-secondary text-left w-full">Username</p>
        <Input icon="material-symbols:person" placeholder="Username" type="text"/>
        <p class="text-text-secondary text-left w-full">Password</p>
        <Input icon="material-symbols:lock" placeholder="Password" type="password"/>
        <button class="button-primary">Login</button>
    </form>
</div>