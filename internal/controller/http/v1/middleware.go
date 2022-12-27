package v1

import (
	"context"
	"fmt"
	"forum/internal/entity"
	"net/http"
	"time"
)

func (h *Handler) CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		foundUser := h.GetExistedSession(w, r)
		if foundUser.Id == 0 {
			h.Errors(w, http.StatusUnauthorized)
			return
		}
		isAuthorized, err := h.Usecases.Users.CheckSession(foundUser)
		if err != nil {
			h.l.WriteLog(fmt.Errorf("v1 - CheckAuth - CheckSession: %w", err))
			h.Errors(w, http.StatusInternalServerError)
			return
		}
		if !isAuthorized {
			h.Errors(w, http.StatusUnauthorized)
			return
		}
		existTTL := foundUser.SessionTTL
		if existTTL.Before(time.Now()) {
			newToken, TTL, err := h.Usecases.Users.GetNewToken()
			if err != nil {
				h.l.WriteLog(fmt.Errorf("v1 - CheckAuth - GetNewToken: %w", err))
				h.Errors(w, http.StatusInternalServerError)
				return
			}
			h.setCookie(w, entity.User{SessionToken: newToken, SessionTTL: TTL})
		}
		err = h.Usecases.Users.UpdateSession(foundUser)
		if err != nil {
			h.l.WriteLog(fmt.Errorf("v1 - CheckAuth - UpdateSession: %w", err))
			h.Errors(w, http.StatusInternalServerError)
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
		isAuthorized, err := h.Usecases.Users.CheckSession(foundUser)
		if err != nil {
			h.l.WriteLog(fmt.Errorf("v1 - AssignStatus - CheckSession: %w", err))
			h.Errors(w, http.StatusInternalServerError)
			return
		}
		if isAuthorized {
			err = h.Usecases.Users.UpdateSession(foundUser)
			if err != nil {
				h.l.WriteLog(fmt.Errorf("v1 - AssignStatus - UpdateSession: %w", err))
				h.Errors(w, http.StatusInternalServerError)
				return
			}
			existTTL := foundUser.SessionTTL
			if existTTL.Before(time.Now()) {
				newToken, TTL, err := h.Usecases.Users.GetNewToken()
				if err != nil {
					h.l.WriteLog(fmt.Errorf("v1 - AssignStatus - GetNewToken: %w", err))
					h.Errors(w, http.StatusInternalServerError)
					return
				}
				h.setCookie(w, entity.User{SessionToken: newToken, SessionTTL: TTL})
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
