<script>
    import Icon from "@iconify/svelte";
    import Input from "$lib/components/Input.svelte";

    let secretKey = $state("");

    function handleSubmit() {
        console.log(secretKey);
    }
</script>

<div class="background-primary-centered">
    <Icon icon="eos-icons:admin" class="text-accent text-left w-full text-6xl mb-4"/>
    <h1 class="text-text-primary text-center font-bold w-full text-2xl mb-4">Admin Dashboard Login</h1>
    <p class="text-text-secondary text-center w-full text-sm mb-4">Enter your credentials to login</p>
    <form action="" class="flex flex-col items-center justify-center bg-background-secondary p-4 rounded-md gap-2 w-1/4">
        <p class="text-text-secondary text-left w-full">Secret key</p>
        <Input icon="material-symbols:lock" placeholder="Secret key" bind:value={secretKey} type="password"/>
        <button type="submit" class="button-primary" onclick={handleSubmit}>Login</button>
    </form>
</div>