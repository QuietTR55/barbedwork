let access_token = $state('');

export function setAccessToken(token: string) {
	access_token = token;
}

export function getAccessToken() {
	return access_token;
}
