<script lang="ts">
    import Icon from "@iconify/svelte";
    let { icon, value = $bindable(""), ...props } = $props();
</script>

<div class="relative w-full">
    <div class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none outline-none">
        <Icon icon={icon} class="text-textSecondary w-5 h-5"/>
    </div>
    <input class="input-primary" bind:value={value} {...props} />
</div>