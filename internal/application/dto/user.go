package dto

import "encoding/json"

type CreateUserReqDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type EventMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}
