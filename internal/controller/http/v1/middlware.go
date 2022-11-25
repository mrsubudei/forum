package v1

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

func (h *Handler) CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		foundUser := h.GetExistedSession(w, r)
		if foundUser.Id == 0 {
			errors.Code = http.StatusUnauthorized
			errors.Message = ErrStatusNotAuthorized
			h.Errors(w, errors)
			return
		}
		isAuthorized, err := h.usecases.Users.CheckSession(foundUser)
		if err != nil {
			log.Println(fmt.Errorf("v1 - CheckAuth - CheckSession: %w", err))
			errors.Code = http.StatusInternalServerError
			errors.Message = ErrInternalServer
			h.Errors(w, errors)
			return
		}
		if !isAuthorized {
			errors.Code = http.StatusUnauthorized
			errors.Message = ErrStatusNotAuthorized
			h.Errors(w, errors)
			return
		}
		err = h.usecases.Users.UpdateSession(foundUser)
		if err != nil {
			log.Println(fmt.Errorf("v1 - CheckAuth - UpdateSession: %w", err))
			errors.Code = http.StatusInternalServerError
			errors.Message = ErrInternalServer
			h.Errors(w, errors)
			return
		}

		content := Content{}
		if foundUser.Id == 1 {
			content.Admin = true
		}
		content.User.Id = foundUser.Id
		content.Authorized = isAuthorized
		content.Unauthorized = !isAuthorized
		ctx := context.Background()
		key := Key("content")

		ctx = context.WithValue(ctx, key, content)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) AssignStatus(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		foundUser := h.GetExistedSession(w, r)
		isAuthorized, err := h.usecases.Users.CheckSession(foundUser)
		if err != nil {
			log.Println(fmt.Errorf("v1 - AssignStatus - CheckSession: %w", err))
			errors.Code = http.StatusInternalServerError
			errors.Message = ErrInternalServer
			h.Errors(w, errors)
			return
		}
		if isAuthorized {
			err = h.usecases.Users.UpdateSession(foundUser)
			if err != nil {
				log.Println(fmt.Errorf("v1 - AssignStatus - UpdateSession: %w", err))
				errors.Code = http.StatusInternalServerError
				errors.Message = ErrInternalServer
				h.Errors(w, errors)
				return
			}
		}
		content := Content{}
		if foundUser.Id == 1 {
			content.Admin = true
		}
		content.User.Id = foundUser.Id
		content.Authorized = isAuthorized
		content.Unauthorized = !isAuthorized
		ctx := context.Background()
		key := Key("content")

		ctx = context.WithValue(ctx, key, content)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
