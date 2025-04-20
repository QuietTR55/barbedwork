<script lang="ts">
	import { themeName, toggleTheme } from '$lib/stores/theme';
    function handleKeyDown(event: KeyboardEvent) {
        if (event.key === 'Enter') {
            console.log('handleKeyDown Enter');
            toggleTheme();
        }
    }

    function handleClick() {
        console.log('handleClick');
        toggleTheme();
    }
</script>

<div class="background-primary-centered">
    <button on:click={handleClick} on:keydown={handleKeyDown} class="button-primary">Toggle Theme</button>
    <p>Current Theme: {$themeName}</p> 
</div>