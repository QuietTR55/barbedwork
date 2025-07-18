<script lang="ts">
	import { onMount } from 'svelte';
	import type { PageData } from './$types';
	import type { Role, Permission, WorkspaceUserRole } from '$lib/models/role';
	import { RoleService, availablePermissions } from '$lib/services/roleService';
	import Input from '$lib/components/Input.svelte';
	import Modal from '$lib/components/ui/Modal.svelte';
	import Icon from '@iconify/svelte';

	let { data }: { data: PageData } = $props();

	// State for role management
	let roles = $state<Role[]>(data.roles);
	let workspaceUserRoles = $state<WorkspaceUserRole[]>(data.workspaceUserRoles);
	let search = $state('');

	// Modal states
	let isCreateRoleModalOpen = $state(false);
	let isEditRoleModalOpen = $state(false);

	// Form states
	let newRoleName = $state('');
	let newRoleDescription = $state('');
	let selectedPermissions = $state<number[]>([]);
	let editingRole: Role | null = $state(null);

	function truncateText(text: string, maxLength = 30) {
		if (text.length <= maxLength) return text;
		return text.substring(0, maxLength) + '...';
	}

	const handleCreateRole = async (event: Event) => {
		event.preventDefault();
		try {
			const permissionObjects = availablePermissions.filter((p) =>
				selectedPermissions.includes(p.id)
			);
			const newRole = await RoleService.createRole({
				name: newRoleName,
				description: newRoleDescription,
				permissions: permissionObjects
			});
			roles = [...roles, newRole];
			resetCreateRoleForm();
			isCreateRoleModalOpen = false;
		} catch (error) {
			console.error('Failed to create role:', error);
		}
	};

	const handleEditRole = async (event: Event) => {
		event.preventDefault();
		if (!editingRole) return;

		try {
			const permissionObjects = availablePermissions.filter((p) =>
				selectedPermissions.includes(p.id)
			);
			const updatedRole = await RoleService.updateRole(editingRole.id, {
				name: newRoleName,
				description: newRoleDescription,
				permissions: permissionObjects
			});

			roles = roles.map((r) => (r.id === editingRole!.id ? updatedRole : r));
			resetCreateRoleForm();
			isEditRoleModalOpen = false;
			editingRole = null;
		} catch (error) {
			console.error('Failed to update role:', error);
		}
	};

	const handleDeleteRole = async (roleId: number) => {
		if (!confirm('Are you sure you want to delete this role?')) return;

		try {
			await RoleService.deleteRole(roleId);
			roles = roles.filter((r) => r.id !== roleId);
			// Also remove any user role assignments
			workspaceUserRoles = workspaceUserRoles.filter((wur) => wur.role_id !== roleId);
		} catch (error) {
			console.error('Failed to delete role:', error);
		}
	};

	const handleAssignRole = async (event: Event) => {
		event.preventDefault();
		alert(
			'Role assignment to users should be done within individual workspaces. Please navigate to a specific workspace to assign roles to users.'
		);
	};

	const handleRemoveUserRole = async (userId: string, roleId: number) => {
		// In admin dashboard, workspace-specific role assignments are not supported
		// This functionality should be implemented within individual workspace management pages
		alert(
			'Role assignment management should be done within individual workspaces. Please navigate to a specific workspace to manage user role assignments.'
		);
	};

	const openEditRoleModal = (role: Role) => {
		editingRole = role;
		newRoleName = role.name;
		newRoleDescription = role.description;
		selectedPermissions = role.permissions.map((p) => p.id);
		isEditRoleModalOpen = true;
	};

	const resetCreateRoleForm = () => {
		newRoleName = '';
		newRoleDescription = '';
		selectedPermissions = [];
	};

	const filteredRoles = $derived(
		roles.filter(
			(role) =>
				role.name.toLowerCase().includes(search.toLowerCase()) ||
				role.description.toLowerCase().includes(search.toLowerCase())
		)
	);
</script>

<div class="flex h-full flex-col gap-4 overflow-y-auto lg:h-max lg:flex-row">
	<!-- Left Panel - Role Management -->
	<div class="bg-background-secondary flex w-full flex-col gap-4 rounded-md p-4 lg:w-1/2">
		<div class="flex flex-col gap-2">
			<div class="flex items-center justify-between">
				<h2 class="text-xl font-bold">Roles</h2>
				<button onclick={() => (isCreateRoleModalOpen = true)} class="button-primary">
					<Icon icon="ic:baseline-add" class="h-5 w-5" />
					<span class="hidden sm:inline">Create Role</span>
				</button>
			</div>
			<Input
				icon="ic:baseline-search"
				placeholder="Search roles..."
				type="text"
				bind:value={search}
				maxlength={50}
			/>
			<p class="text-text-secondary text-sm">Showing {filteredRoles.length} roles</p>
		</div>

		<div class="flex max-h-96 flex-col gap-2 overflow-y-auto">
			{#each filteredRoles as role (role.id)}
				<div class="bg-background-tertiary rounded-md p-3">
					<div class="flex items-center justify-between">
						<div class="flex-1">
							<h3 class="font-semibold">{role.name}</h3>
							<p class="text-text-secondary text-sm">{truncateText(role.description, 40)}</p>
							<p class="text-text-secondary mt-1 text-xs">
								{role.permissions.length} permission{role.permissions.length !== 1 ? 's' : ''}
							</p>
						</div>
						<div class="flex gap-2">
							<button onclick={() => openEditRoleModal(role)} class="button-secondary">
								<Icon icon="ic:baseline-edit" class="h-4 w-4" />
							</button>
							<button onclick={() => handleDeleteRole(role.id)} class="button-important">
								<Icon icon="ic:baseline-delete" class="h-4 w-4" />
							</button>
						</div>
					</div>

					<!-- Permissions preview -->
					<div class="mt-2 flex flex-wrap gap-1">
						{#each role.permissions.slice(0, 3) as permission}
							<span class="bg-accent text-text-button rounded px-2 py-1 text-xs">
								{permission.name.replace('workspace:', '')}
							</span>
						{/each}
						{#if role.permissions.length > 3}
							<span class="text-text-secondary px-2 py-1 text-xs">
								+{role.permissions.length - 3} more
							</span>
						{/if}
					</div>
				</div>
			{/each}
		</div>
	</div>

	<!-- Right Panel - Global Role Information -->
	<div class="bg-background-secondary flex w-full flex-col gap-4 rounded-md p-4 lg:w-1/2">
		<div class="flex items-center justify-between">
			<h2 class="text-xl font-bold">Global Role Management</h2>
		</div>

		<div class="bg-background-primary rounded-md p-4">
			<div class="text-center">
				<Icon icon="ic:baseline-info" class="mx-auto mb-2 h-12 w-12 text-blue-500" />
				<h3 class="mb-2 font-semibold">Admin Dashboard - Global Roles</h3>
				<p class="text-text-secondary mb-4 text-sm">
					This interface allows you to manage global role definitions and their permissions. To
					assign roles to specific users, please navigate to individual workspace management pages.
				</p>
				<div class="text-text-secondary space-y-1 text-xs">
					<p><strong>Current functionality:</strong></p>
					<ul class="list-inside list-disc space-y-1 text-left">
						<li>Create and edit global role definitions</li>
						<li>Manage permissions for each role</li>
						<li>Delete unused roles</li>
					</ul>
					<p class="mt-3"><strong>For user assignments:</strong></p>
					<p>
						Navigate to <code class="bg-background-tertiary rounded px-1"
							>/workspaces/[workspace-id]</code
						> to manage user role assignments within specific workspaces.
					</p>
				</div>
			</div>
		</div>
	</div>
</div>

<!-- Create Role Modal -->
<Modal
	title="Create New Role"
	isOpen={isCreateRoleModalOpen}
	onClose={() => {
		isCreateRoleModalOpen = false;
		resetCreateRoleForm();
	}}
>
	<form onsubmit={handleCreateRole} class="flex flex-col gap-4">
		<Input
			icon="ic:baseline-label"
			placeholder="Role Name"
			type="text"
			bind:value={newRoleName}
			maxlength={50}
			required
		/>
		<Input
			icon="ic:baseline-description"
			placeholder="Role Description"
			type="text"
			bind:value={newRoleDescription}
			maxlength={200}
			required
		/>

		<div>
			<h3 class="mb-2 font-semibold">Permissions</h3>
			<div class="border-border-primary max-h-48 overflow-y-auto rounded border p-2">
				{#each availablePermissions as permission}
					<label class="hover:bg-background-tertiary flex items-center gap-2 rounded p-2">
						<input
							type="checkbox"
							bind:group={selectedPermissions}
							value={permission.id}
							class="accent-accent"
						/>
						<div>
							<p class="font-medium">{permission.name}</p>
							<p class="text-text-secondary text-sm">{permission.description}</p>
						</div>
					</label>
				{/each}
			</div>
		</div>

		<div class="flex justify-end gap-2">
			<button
				type="button"
				onclick={() => {
					isCreateRoleModalOpen = false;
					resetCreateRoleForm();
				}}
				class="button-secondary"
			>
				Cancel
			</button>
			<button type="submit" class="button-primary"> Create Role </button>
		</div>
	</form>
</Modal>

<!-- Edit Role Modal -->
<Modal
	title="Edit Role"
	isOpen={isEditRoleModalOpen}
	onClose={() => {
		isEditRoleModalOpen = false;
		editingRole = null;
		resetCreateRoleForm();
	}}
>
	<form onsubmit={handleEditRole} class="flex flex-col gap-4">
		<Input
			icon="ic:baseline-label"
			placeholder="Role Name"
			type="text"
			bind:value={newRoleName}
			maxlength={50}
			required
		/>
		<Input
			icon="ic:baseline-description"
			placeholder="Role Description"
			type="text"
			bind:value={newRoleDescription}
			maxlength={200}
			required
		/>

		<div>
			<h3 class="mb-2 font-semibold">Permissions</h3>
			<div class="border-border-primary max-h-48 overflow-y-auto rounded border p-2">
				{#each availablePermissions as permission}
					<label class="hover:bg-background-tertiary flex items-center gap-2 rounded p-2">
						<input
							type="checkbox"
							bind:group={selectedPermissions}
							value={permission.id}
							class="accent-accent"
						/>
						<div>
							<p class="font-medium">{permission.name}</p>
							<p class="text-text-secondary text-sm">{permission.description}</p>
						</div>
					</label>
				{/each}
			</div>
		</div>

		<div class="flex justify-end gap-2">
			<button
				type="button"
				onclick={() => {
					isEditRoleModalOpen = false;
					editingRole = null;
					resetCreateRoleForm();
				}}
				class="button-secondary"
			>
				Cancel
			</button>
			<button type="submit" class="button-primary"> Update Role </button>
		</div>
	</form>
</Modal>
