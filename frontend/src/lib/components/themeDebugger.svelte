<script lang="ts">
    import { themeName, toggleTheme } from '$lib/stores/theme';
</script>

<button onclick={() => toggleTheme()} class="bg-background-primary text-text-primary rounded-md border-2 p-2 fixed bottom-0 right-0">
    <p>Current Theme: {$themeName}</p>
</button>
