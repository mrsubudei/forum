package v1

import (
	"forum/internal/usecase"
	"net/http"
	"text/template"
)

type Handler struct {
	usecases *usecase.UseCases
}

func NewHandler(usecases *usecase.UseCases) *Handler {
	return &Handler{
		usecases: usecases,
	}
}

func (h *Handler) Errors(w http.ResponseWriter, errors ErrMessage) {
	html, err := template.ParseFiles("templates/errors.html")
	if err != nil {
		http.Error(w, ErrInternalServer, http.StatusInternalServerError)
		return
	}

	html.Execute(w, errors)

}
