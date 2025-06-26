package role

type RoleCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
}

type RoleUpdateRequest struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
}

type RoleResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}
