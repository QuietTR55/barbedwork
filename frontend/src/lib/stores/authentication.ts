import { writable } from "svelte/store";

let access_token = writable<string>("");

export function setAccessToken(token: string) {
	access_token.set(token);
}

export function getAccessToken() {
	return access_token.subscribe(value => value);
}
