package v1

import (
	"fmt"
	"html/template"
	"net/http"

	"forum/internal/config"
	"forum/internal/usecase"
	"forum/pkg/logger"
)

type Handler struct {
	usecases *usecase.UseCases
	Cfg      config.Config
	l        *logger.Logger
}

func NewHandler(usecases *usecase.UseCases, cfg config.Config, logger *logger.Logger) *Handler {
	return &Handler{
		usecases: usecases,
		Cfg:      cfg,
		l:        logger,
	}
}

func (h *Handler) ParseAndExecute(w http.ResponseWriter, content Content, path string) error {
	html, err := template.ParseFiles(path)
	if err != nil {
		h.l.WriteLog(fmt.Errorf("parseFiles: %w", err))
		h.Errors(w, http.StatusInternalServerError)
		return err
	}

	err = html.Execute(w, content)
	if err != nil {
		h.l.WriteLog(fmt.Errorf("execute: %w", err))
		h.Errors(w, http.StatusInternalServerError)
		return err
	}
	return nil
}

func (h *Handler) Errors(w http.ResponseWriter, status int) {
	errors := ErrMessage{}
	switch status {
	case http.StatusBadRequest:
		errors.Code = http.StatusBadRequest
		errors.Message = ErrBadRequest
	case http.StatusUnauthorized:
		errors.Code = http.StatusUnauthorized
		errors.Message = ErrStatusNotAuthorized
	case http.StatusForbidden:
		errors.Code = http.StatusForbidden
		errors.Message = ErrLowAccessLevel
	case http.StatusNotFound:
		errors.Code = http.StatusNotFound
		errors.Message = ErrPageNotFound
	case http.StatusMethodNotAllowed:
		errors.Code = http.StatusMethodNotAllowed
		errors.Message = ErrMethodNotAllowed
	case http.StatusNotAcceptable:
		errors.Code = http.StatusNotAcceptable
		errors.Message = UserNotExist
	case http.StatusInternalServerError:
		errors.Code = http.StatusInternalServerError
		errors.Message = ErrInternalServer
	}
	html, err := template.ParseFiles("templates/errors.html")
	if err != nil {
		h.l.WriteLog(fmt.Errorf("v1 - Errors - ParseFiles: %w", err))
		http.Error(w, ErrInternalServer, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(errors.Code)
	html.Execute(w, errors)
}
