package dto

type SigninRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SigninResponse struct {
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	// RuleID   uint   `json:"rule_id"`
	// StatusID uint   `json:"status_id"`
	Avatar string `json:"avatar"`
	Rule   RuleResponse
	Status StatusResponse
}

type RuleResponse struct {
	ID       uint   `json:"id"`
	RuleName string `'json:"rule_name"`
}
type StatusResponse struct {
	ID         uint   `json:"id"`
	StatusName string `'json:"status_name"`
}
