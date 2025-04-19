package dto

import "gorm.io/datatypes"

type EventReqDTO struct {
	Type    string         `json:"type"`
	Payload datatypes.JSON `json:"payload"`
}
