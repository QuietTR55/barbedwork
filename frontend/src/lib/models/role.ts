export type Permission = {
  id: number; // SERIAL PRIMARY KEY
  name: string;
  description: string;
};

export type Role = {
  id: number; // SERIAL PRIMARY KEY  
  name: string;
  description: string;
  permissions: Permission[];
  created_at: string;
  updated_at: string;
};

export type WorkspaceUserRole = {
  workspace_id: string; // UUID
  user_id: string; // UUID
  role_id: number; // INT
  assigned_at: string; // TIMESTAMP
};