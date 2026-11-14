package dto

type TenantRequest struct {
	Name      string `json:"name" binding:"required"`
	Subdomain string `json:"subdomain" binding:"required"`
	Domain    string `json:"domain"`
	Status    string `json:"status"`
}
