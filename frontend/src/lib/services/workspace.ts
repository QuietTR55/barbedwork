export const getWorkspaceAsAdmin = async (workspaceId: string): Promise<any> => {
	const endPoint = '/api/admin/workspaces/' + workspaceId;
	const response = await fetch(endPoint, {
		method: 'GET',
		headers: {
			'Content-Type': 'application/json'
		},
		credentials: 'include'
	});

	if (!response.ok) {
		throw new Error('Failed to fetch workspace');
	}

	const workspace = await response.json();
	return workspace;
};
