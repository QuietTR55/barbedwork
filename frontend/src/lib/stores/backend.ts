import { get, writable } from "svelte/store";
import { browser } from "$app/environment";

const BACKEND_URL_STORAGE_KEY = "backendUrl";

let backendUrl = writable<string>("http://localhost:3000");

// Check for stored backend URL only in the browser
if (browser) {
  const storedBackendUrl = localStorage.getItem(BACKEND_URL_STORAGE_KEY);
  if (storedBackendUrl) {
    backendUrl.set(storedBackendUrl);
  }
}

const setBackendUrl = (url: string) => {
  backendUrl.set(url);
  if (browser) {
    localStorage.setItem(BACKEND_URL_STORAGE_KEY, url);
  }
};

export { backendUrl, setBackendUrl };