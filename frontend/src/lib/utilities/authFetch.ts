import { getAccessToken, setAccessToken } from '$lib/stores/authentication';

export const authFetch = async (url: string, options: RequestInit = {}): Promise<Response> => {
  const accessToken = getAccessToken();
  if (!accessToken) {
    throw new Error('No access token found');
  }

  const headers: Record<string, string> = {
    ...(options.headers as Record<string, string> || {}),
    Authorization: `Bearer ${accessToken}`
  };

  if (!(options.body instanceof FormData)) {
    headers['Content-Type'] = 'application/json';
  }

  const response = await fetch(url, {
    ...options,
    headers,
    credentials: 'include'
  });

  const newToken = response.headers.get('X-New-Access-Token');
  if (newToken) {
    setAccessToken(newToken);
  }

  return response;
};
