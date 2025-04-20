import { browser } from '$app/environment';
import { writable, get } from 'svelte/store';

const THEME_NAME_STORAGE_KEY = 'theme_name';

function getInitialTheme(): 'light' | 'dark' {
	if (!browser) return 'dark';
	const stored = localStorage.getItem(THEME_NAME_STORAGE_KEY);
	return stored === 'light' || stored === 'dark' ? stored : 'dark';
}

const themeName = writable<'light' | 'dark'>(getInitialTheme());

if (browser) {
	themeName.subscribe(($themeName) => {
		localStorage.setItem(THEME_NAME_STORAGE_KEY, $themeName);
		document.documentElement.setAttribute('data-theme', $themeName);
		document.documentElement.classList.toggle('light', $themeName === 'light');
		document.documentElement.classList.toggle('dark', $themeName === 'dark');
	});
}

function toggleTheme() {
	console.log('toggleTheme');
	themeName.update((current) => (current === 'light' ? 'dark' : 'light'));
}

export { themeName, toggleTheme };
