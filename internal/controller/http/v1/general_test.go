package v1_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"forum/internal/config"
	v1 "forum/internal/controller/http/v1"
	"forum/internal/usecase"
	mu "forum/internal/usecase/mock"
	"forum/pkg/logger"
)

func TestIndexHandler(t *testing.T) {
	cfg, err := config.LoadConfig("../../../../config.json")
	if err != nil {
		log.Fatal(err)
	}
	l := logger.New()
	mockUsersUseCase := mu.NewUsersMockUseCase()
	mockPostsUseCase := mu.NewPostsMockUseCase()
	mockCommentsUseCase := mu.NewCommentsMockUseCase()
	usecases := usecase.NewUseCases(mockPostsUseCase, mockUsersUseCase, mockCommentsUseCase)
	handler := v1.NewHandler(usecases, cfg, l)
	handler.RegisterRoutes(handler.Mux)

	t.Run("OK", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("want: %v, got: %v", http.StatusOK, rec.Code)
		}
	})

	t.Run("err wrong method", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/", nil)
		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("want: %v, got: %v", http.StatusMethodNotAllowed, rec.Code)
		}
	})
}
