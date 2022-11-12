package v1

import (
	"forum/internal/usecase"
)

type Handler struct {
	services *usecase.UseCases
}

func NewHandler(services *usecase.UseCases) *Handler {
	return &Handler{
		services: services,
	}
}
