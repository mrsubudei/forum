package v1

import (
	"forum/internal/usecase"
)

type Handler struct {
	usecases *usecase.UseCases
}

func NewHandler(services *usecase.UseCases) *Handler {
	return &Handler{
		usecases: services,
	}
}
