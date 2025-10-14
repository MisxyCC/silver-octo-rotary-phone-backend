package models

type ApprovalRequest struct {
	User    string `json:"user" binding:"required"`
	Amount  int    `json:"amount" binding:"required"`
	Details string `json:"details" binding:"required"`
}
