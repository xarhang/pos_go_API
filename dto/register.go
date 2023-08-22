package dto

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Fullname string `json:"fullname" binding:"required"`
	Avatar   string `json:"avatar" binding:"required"`
	Status   int    `json:"status" binding:"required"`
	Rule     int    `json:"rule" binding:"required"`
}

type RegisterResponse struct {
	ID       int    `json:"id" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Fullname string `json:"fullname" binding:"required"`
	Avatar   string `json:"avatar" binding:"required"`
	Status   int    `json:"status" binding:"required"`
	Rule     int    `json:"rule" binding:"required"`
}
