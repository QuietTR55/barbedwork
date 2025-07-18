package handlers

import (
	"backend/internal/models"
	"backend/internal/services"
	"backend/pkg/middleware"
	"backend/pkg/ratelimiter"
	"backend/pkg/utilities"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type RoleHandler struct {
	roleService       *services.RoleService
	store             utilities.SessionStore
	limiter           ratelimiter.RateLimiter
	permissionChecker *utilities.PermissionChecker
}

func NewRoleHandler(roleService *services.RoleService, store utilities.SessionStore, limiter ratelimiter.RateLimiter, permissionChecker *utilities.PermissionChecker) *RoleHandler {
	return &RoleHandler{
		roleService:       roleService,
		store:             store,
		limiter:           limiter,
		permissionChecker: permissionChecker,
	}
}

func (h *RoleHandler) RegisterRoutes(router *http.ServeMux) {
	log.Println("RoleHandler: Registering routes")
	
	// Admin middleware stack for role management - allow multiple permissions
	adminStack := []middleware.Middleware{
		middleware.TokenAuthMiddleware(h.store),
		middleware.RateLimitMiddleware(h.limiter, time.Minute, "roles"),
		middleware.PermissionMiddleware(h.permissionChecker, "workspace:manage-roles", "workspace:manage-workspace"),
	}

	// User middleware stack for viewing roles
	userStack := []middleware.Middleware{
		middleware.TokenAuthMiddleware(h.store),
		middleware.RateLimitMiddleware(h.limiter, time.Minute, "roles_view"),
	}

	log.Println("RoleHandler: Registering /api/permissions route")
	// Permission routes
	router.Handle("/api/permissions", middleware.Chain(
		http.HandlerFunc(h.GetAllPermissions),
		userStack...,
	))

	log.Println("RoleHandler: Registering /api/roles route")
	// Role management routes
	router.Handle("/api/roles", middleware.Chain(
		http.HandlerFunc(h.handleRoles),
		adminStack...,
	))

	log.Println("RoleHandler: Registering /api/roles/ route")
	router.Handle("/api/roles/", middleware.Chain(
		http.HandlerFunc(h.handleRoleByID),
		adminStack...,
	))

	log.Println("RoleHandler: Registering /api/workspaces/ route")
	// Workspace user role assignment routes
	router.Handle("/api/workspaces/", middleware.Chain(
		http.HandlerFunc(h.handleWorkspaceRoles),
		adminStack...,
	))
	
	log.Println("RoleHandler: All routes registered successfully")
}

// GetAllPermissions retrieves all available permissions
func (h *RoleHandler) GetAllPermissions(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetAllPermissions: Request received - Method: %s, URL: %s", r.Method, r.URL.Path)
	
	if r.Method != http.MethodGet {
		log.Printf("GetAllPermissions: Method not allowed - %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Println("GetAllPermissions: Calling roleService.GetAllPermissions")
	permissions, err := h.roleService.GetAllPermissions(r.Context())
	if err != nil {
		log.Printf("GetAllPermissions: Error from roleService.GetAllPermissions: %v", err)
		http.Error(w, "Failed to get permissions", http.StatusInternalServerError)
		return
	}

	log.Printf("GetAllPermissions: Successfully retrieved %d permissions", len(permissions))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(permissions)
	log.Println("GetAllPermissions: Response sent successfully")
}

// handleRoles handles /api/roles endpoint
func (h *RoleHandler) handleRoles(w http.ResponseWriter, r *http.Request) {
	log.Printf("handleRoles: Request received - Method: %s, URL: %s", r.Method, r.URL.Path)
	
	switch r.Method {
	case http.MethodGet:
		log.Println("handleRoles: Routing to GetAllRoles")
		h.GetAllRoles(w, r)
	case http.MethodPost:
		log.Println("handleRoles: Routing to CreateRole")
		h.CreateRole(w, r)
	default:
		log.Printf("handleRoles: Method not allowed - %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetAllRoles retrieves all roles with their permissions
func (h *RoleHandler) GetAllRoles(w http.ResponseWriter, r *http.Request) {
	log.Printf("GetAllRoles: Request received - Method: %s, URL: %s", r.Method, r.URL.Path)
	
	log.Println("GetAllRoles: Calling roleService.GetAllRoles")
	roles, err := h.roleService.GetAllRoles(r.Context())
	if err != nil {
		log.Printf("GetAllRoles: Error from roleService.GetAllRoles: %v", err)
		http.Error(w, "Failed to get roles", http.StatusInternalServerError)
		return
	}

	log.Printf("GetAllRoles: Successfully retrieved %d roles", len(roles))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roles)
	log.Println("GetAllRoles: Response sent successfully")
}

// CreateRole creates a new role
func (h *RoleHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	log.Printf("CreateRole: Request received - Method: %s, URL: %s", r.Method, r.URL.Path)
	
	var req models.CreateRoleRequest
	log.Println("CreateRole: Decoding request body")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("CreateRole: Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	var descStr string
	if req.Description != nil {
		descStr = *req.Description
	} else {
		descStr = "<nil>"
	}
	log.Printf("CreateRole: Request decoded - Name: %s, Description: %s, PermissionIDs: %v", req.Name, descStr, req.PermissionIDs)

	log.Println("CreateRole: Calling roleService.CreateRole")
	role, err := h.roleService.CreateRole(r.Context(), req)
	if err != nil {
		log.Printf("CreateRole: Error from roleService.CreateRole: %v", err)
		switch err {
		case services.ErrInvalidRoleName:
			log.Println("CreateRole: Invalid role name error")
			http.Error(w, err.Error(), http.StatusBadRequest)
		case services.ErrRoleNameExists:
			log.Println("CreateRole: Role name exists error")
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			log.Println("CreateRole: Internal server error")
			http.Error(w, "Failed to create role", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("CreateRole: Role created successfully - ID: %d, Name: %s", role.ID, role.Name)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(role)
	log.Println("CreateRole: Response sent successfully")
}

// handleRoleByID handles /api/roles/{id} endpoint
func (h *RoleHandler) handleRoleByID(w http.ResponseWriter, r *http.Request) {
	log.Printf("handleRoleByID: Request received - Method: %s, URL: %s", r.Method, r.URL.Path)
	
	// Extract role ID from path
	path := strings.TrimPrefix(r.URL.Path, "/api/roles/")
	log.Printf("handleRoleByID: Extracted path: %s", path)
	
	roleID, err := strconv.Atoi(path)
	if err != nil {
		log.Printf("handleRoleByID: Error converting roleID to int: %v", err)
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}
	
	log.Printf("handleRoleByID: Parsed roleID: %d", roleID)

	switch r.Method {
	case http.MethodGet:
		log.Println("handleRoleByID: Routing to GetRoleByID")
		h.GetRoleByID(w, r, roleID)
	case http.MethodPut:
		log.Println("handleRoleByID: Routing to UpdateRole")
		h.UpdateRole(w, r, roleID)
	case http.MethodDelete:
		log.Println("handleRoleByID: Routing to DeleteRole")
		h.DeleteRole(w, r, roleID)
	default:
		log.Printf("handleRoleByID: Method not allowed - %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetRoleByID retrieves a specific role by ID
func (h *RoleHandler) GetRoleByID(w http.ResponseWriter, r *http.Request, roleID int) {
	log.Printf("GetRoleByID: Request received for roleID: %d", roleID)
	
	log.Printf("GetRoleByID: Calling roleService.GetRoleByID for ID: %d", roleID)
	role, err := h.roleService.GetRoleByID(r.Context(), roleID)
	if err != nil {
		log.Printf("GetRoleByID: Error from roleService.GetRoleByID: %v", err)
		switch err {
		case services.ErrRoleNotFound:
			log.Printf("GetRoleByID: Role not found for ID: %d", roleID)
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			log.Printf("GetRoleByID: Internal server error for ID: %d", roleID)
			http.Error(w, "Failed to get role", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("GetRoleByID: Successfully retrieved role - ID: %d, Name: %s", role.ID, role.Name)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(role)
	log.Println("GetRoleByID: Response sent successfully")
}

// UpdateRole updates an existing role
func (h *RoleHandler) UpdateRole(w http.ResponseWriter, r *http.Request, roleID int) {
	log.Printf("UpdateRole: Request received for roleID: %d", roleID)
	
	var req models.UpdateRoleRequest
	log.Println("UpdateRole: Decoding request body")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("UpdateRole: Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	var descStr string
	if req.Description != nil {
		descStr = *req.Description
	} else {
		descStr = "<nil>"
	}
	log.Printf("UpdateRole: Request decoded - Name: %s, Description: %s, PermissionIDs: %v", req.Name, descStr, req.PermissionIDs)

	log.Printf("UpdateRole: Calling roleService.UpdateRole for ID: %d", roleID)
	role, err := h.roleService.UpdateRole(r.Context(), roleID, req)
	if err != nil {
		log.Printf("UpdateRole: Error from roleService.UpdateRole: %v", err)
		switch err {
		case services.ErrRoleNotFound:
			log.Printf("UpdateRole: Role not found for ID: %d", roleID)
			http.Error(w, err.Error(), http.StatusNotFound)
		case services.ErrInvalidRoleName:
			log.Println("UpdateRole: Invalid role name error")
			http.Error(w, err.Error(), http.StatusBadRequest)
		case services.ErrRoleNameExists:
			log.Println("UpdateRole: Role name exists error")
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			log.Printf("UpdateRole: Internal server error for ID: %d", roleID)
			http.Error(w, "Failed to update role", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("UpdateRole: Role updated successfully - ID: %d, Name: %s", role.ID, role.Name)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(role)
	log.Println("UpdateRole: Response sent successfully")
}

// DeleteRole deletes a role
func (h *RoleHandler) DeleteRole(w http.ResponseWriter, r *http.Request, roleID int) {
	log.Printf("DeleteRole: Request received for roleID: %d", roleID)
	
	log.Printf("DeleteRole: Calling roleService.DeleteRole for ID: %d", roleID)
	err := h.roleService.DeleteRole(r.Context(), roleID)
	if err != nil {
		log.Printf("DeleteRole: Error from roleService.DeleteRole: %v", err)
		switch err {
		case services.ErrRoleNotFound:
			log.Printf("DeleteRole: Role not found for ID: %d", roleID)
			http.Error(w, err.Error(), http.StatusNotFound)
		default:
			log.Printf("DeleteRole: Internal server error for ID: %d", roleID)
			http.Error(w, "Failed to delete role", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("DeleteRole: Role deleted successfully - ID: %d", roleID)
	w.WriteHeader(http.StatusNoContent)
	log.Println("DeleteRole: Response sent successfully")
}

// handleWorkspaceRoles handles workspace role assignment endpoints
func (h *RoleHandler) handleWorkspaceRoles(w http.ResponseWriter, r *http.Request) {
	log.Printf("handleWorkspaceRoles: Request received - Method: %s, URL: %s", r.Method, r.URL.Path)
	
	// Parse path: /api/workspaces/{workspace_id}/user-roles
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/workspaces/"), "/")
	log.Printf("handleWorkspaceRoles: Path parts: %v", pathParts)
	
	if len(pathParts) < 2 {
		log.Printf("handleWorkspaceRoles: Invalid path - insufficient parts: %v", pathParts)
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	workspaceID := pathParts[0]
	endpoint := pathParts[1]
	log.Printf("handleWorkspaceRoles: Parsed - WorkspaceID: %s, Endpoint: %s", workspaceID, endpoint)

	if endpoint != "user-roles" {
		log.Printf("handleWorkspaceRoles: Invalid endpoint: %s", endpoint)
		http.Error(w, "Invalid endpoint", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		log.Println("handleWorkspaceRoles: Routing to GetWorkspaceUserRoles")
		h.GetWorkspaceUserRoles(w, r, workspaceID)
	case http.MethodPost:
		log.Println("handleWorkspaceRoles: Routing to AssignRoleToUser")
		h.AssignRoleToUser(w, r, workspaceID)
	case http.MethodDelete:
		log.Println("handleWorkspaceRoles: Routing to RemoveRoleFromUser")
		h.RemoveRoleFromUser(w, r, workspaceID)
	default:
		log.Printf("handleWorkspaceRoles: Method not allowed - %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetWorkspaceUserRoles retrieves all user role assignments for a workspace
func (h *RoleHandler) GetWorkspaceUserRoles(w http.ResponseWriter, r *http.Request, workspaceID string) {
	log.Printf("GetWorkspaceUserRoles: Request received for workspaceID: %s", workspaceID)
	
	log.Printf("GetWorkspaceUserRoles: Calling roleService.GetWorkspaceUserRoles for workspaceID: %s", workspaceID)
	userRoles, err := h.roleService.GetWorkspaceUserRoles(r.Context(), workspaceID)
	if err != nil {
		log.Printf("GetWorkspaceUserRoles: Error from roleService.GetWorkspaceUserRoles: %v", err)
		http.Error(w, "Failed to get workspace user roles", http.StatusInternalServerError)
		return
	}

	log.Printf("GetWorkspaceUserRoles: Successfully retrieved %d user roles for workspace: %s", len(userRoles), workspaceID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userRoles)
	log.Println("GetWorkspaceUserRoles: Response sent successfully")
}

// AssignRoleToUser assigns a role to a user in a workspace
func (h *RoleHandler) AssignRoleToUser(w http.ResponseWriter, r *http.Request, workspaceID string) {
	log.Printf("AssignRoleToUser: Request received for workspaceID: %s", workspaceID)
	
	var req models.AssignRoleRequest
	log.Println("AssignRoleToUser: Decoding request body")
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("AssignRoleToUser: Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	log.Printf("AssignRoleToUser: Request decoded - UserID: %s, RoleID: %d, WorkspaceID: %s", req.UserID, req.RoleID, workspaceID)

	log.Printf("AssignRoleToUser: Calling roleService.AssignRoleToUser")
	userRole, err := h.roleService.AssignRoleToUser(r.Context(), workspaceID, req.UserID, req.RoleID)
	if err != nil {
		log.Printf("AssignRoleToUser: Error from roleService.AssignRoleToUser: %v", err)
		switch err {
		case services.ErrInvalidUserID, services.ErrInvalidRoleID:
			log.Printf("AssignRoleToUser: Invalid ID error: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		case services.ErrRoleNotFound:
			log.Printf("AssignRoleToUser: Role not found error: %v", err)
			http.Error(w, err.Error(), http.StatusNotFound)
		case services.ErrRoleAssignmentExists:
			log.Printf("AssignRoleToUser: Role assignment exists error: %v", err)
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			log.Printf("AssignRoleToUser: Internal server error: %v", err)
			http.Error(w, "Failed to assign role", http.StatusInternalServerError)
		}
		return
	}

	log.Printf("AssignRoleToUser: Role assigned successfully - UserID: %s, RoleID: %d, WorkspaceID: %s", req.UserID, req.RoleID, workspaceID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userRole)
	log.Println("AssignRoleToUser: Response sent successfully")
}

// RemoveRoleFromUser removes a role assignment from a user
func (h *RoleHandler) RemoveRoleFromUser(w http.ResponseWriter, r *http.Request, workspaceID string) {
	log.Printf("RemoveRoleFromUser: Request received for workspaceID: %s", workspaceID)
	
	// Parse query parameters for user_id and role_id
	userID := r.URL.Query().Get("user_id")
	roleIDStr := r.URL.Query().Get("role_id")
	log.Printf("RemoveRoleFromUser: Query parameters - userID: %s, roleIDStr: %s", userID, roleIDStr)

	if userID == "" || roleIDStr == "" {
		log.Println("RemoveRoleFromUser: Missing required query parameters")
		http.Error(w, "user_id and role_id query parameters are required", http.StatusBadRequest)
		return
	}

	log.Printf("RemoveRoleFromUser: Converting roleIDStr to int: %s", roleIDStr)
	roleID, err := strconv.Atoi(roleIDStr)
	if err != nil {
		log.Printf("RemoveRoleFromUser: Error converting roleID to int: %v", err)
		http.Error(w, "Invalid role_id", http.StatusBadRequest)
		return
	}
	
	log.Printf("RemoveRoleFromUser: Parsed values - UserID: %s, RoleID: %d, WorkspaceID: %s", userID, roleID, workspaceID)

	log.Println("RemoveRoleFromUser: Calling roleService.RemoveRoleFromUser")
	err = h.roleService.RemoveRoleFromUser(r.Context(), workspaceID, userID, roleID)
	if err != nil {
		log.Printf("RemoveRoleFromUser: Error from roleService.RemoveRoleFromUser: %v", err)
		switch err {
		case services.ErrInvalidUserID, services.ErrInvalidRoleID:
			log.Printf("RemoveRoleFromUser: Invalid ID error: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		default:
			if err.Error() == "role assignment not found" {
				log.Printf("RemoveRoleFromUser: Role assignment not found: %v", err)
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				log.Printf("RemoveRoleFromUser: Internal server error: %v", err)
				http.Error(w, "Failed to remove role assignment", http.StatusInternalServerError)
			}
		}
		return
	}

	log.Printf("RemoveRoleFromUser: Role removed successfully - UserID: %s, RoleID: %d, WorkspaceID: %s", userID, roleID, workspaceID)
	w.WriteHeader(http.StatusNoContent)
	log.Println("RemoveRoleFromUser: Response sent successfully")
}
