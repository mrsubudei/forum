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

func (h *Handler) ParseAndExecute(w http.ResponseWriter, content Content, path string) error {
	html, err := template.ParseFiles(path)
	if err != nil {
		log.Println(fmt.Errorf("v1 - ParseAndExecute - ParseFiles: %w", err))
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return err
	}

	err = html.Execute(w, content)
	if err != nil {
		log.Println(fmt.Errorf("v1 - ParseAndExecute - Execute: %w", err))
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
		h.Errors(w, errors)
		return err
	}
	return nil
}

func (h *Handler) Errors(w http.ResponseWriter, errors ErrMessage) {
	html, err := template.ParseFiles("templates/errors.html")
	if err != nil {
		log.Println(fmt.Errorf("v1 - Errors - ParseFiles: %w", err))
		http.Error(w, ErrInternalServer, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(errors.Code)
	html.Execute(w, errors)
}
