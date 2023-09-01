package dto

type VPermissionResponse struct {
	// ID         uint   `json:"id"`
	// Ptype      string `json:"ptype"`
	Path       string `json:"path"`
	Permission string `json:"permission"`
	Rules      string `json:"rules"`
}
type VPermission struct {
	// ID         uint   `json:"id"`
	// Ptype      string `json:"ptype"`
	Path       string `json:"path"`
	Permission string `json:"permission"`
	Rules      string `json:"rules"`
}
type VPermissionrequest struct {
	Ptype      string `json:"ptype" binding:"required"`
	Rules      string `json:"rules" binding:"required"`
	Path       string `json:"path" binding:"required"`
	Permission string `json:"permission" binding:"required"`
}
