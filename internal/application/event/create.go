package event

import "poc-event-source/internal/application/dto"

type MainHandler interface {
	Handler(dto dto.EventReqDTO) error
}
