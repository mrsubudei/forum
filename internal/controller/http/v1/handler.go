package v1

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"forum/internal/config"
	"forum/internal/usecase"
)

type Handler struct {
	usecases *usecase.UseCases
	Cfg      config.Config
}

func NewHandler(usecases *usecase.UseCases, cfg config.Config) *Handler {
	return &Handler{
		usecases: usecases,
		Cfg:      cfg,
	}
}

func (h *Handler) Errors(w http.ResponseWriter, errors ErrMessage) {
	html, err := template.ParseFiles("templates/errors.html")
	if err != nil {
		log.Println(fmt.Errorf("v1 - Errors - ParseFiles: %w", err))
		http.Error(w, ErrInternalServer, http.StatusInternalServerError)
		return
	}

	html.Execute(w, errors)
}
