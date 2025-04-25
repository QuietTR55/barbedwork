import { writable, get } from 'svelte/store';
import { browser } from '$app/environment';

const tokenStorageKey = 'access_token';

const initialTokenValue = browser ? (localStorage.getItem(tokenStorageKey) ?? '') : '';
let access_token = writable<string>(initialTokenValue);

export function setAccessToken(token: string) {
	if (!browser) return;
	access_token.set(token);
	localStorage.setItem(tokenStorageKey, token);
}

export function getAccessToken(): string {
	if (!browser) return '';

	const token = get(access_token);
	return token;
}

export function clearAccessToken() {
	if (!browser) return;

	access_token.set('');
	localStorage.removeItem(tokenStorageKey);
}

if (browser) {
	access_token.subscribe((value) => {
		if (value !== localStorage.getItem(tokenStorageKey)) {
			if (value) {
				localStorage.setItem(tokenStorageKey, value);
			} else {
				localStorage.removeItem(tokenStorageKey);
			}
		}
	});
}
