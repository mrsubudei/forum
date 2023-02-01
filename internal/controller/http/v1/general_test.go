package v1_test

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"forum/internal/config"
	v1 "forum/internal/controller/http/v1"
	"forum/internal/entity"
	"forum/internal/usecase"
	mu "forum/internal/usecase/mock"
	"forum/pkg/logger"
)

func setup() *v1.Handler {
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

	return handler
}

func CreateMultipartForm(t *testing.T, path string,
	content string, fieldName string) (*bytes.Buffer, *multipart.Writer) {

	body := new(bytes.Buffer)
	mw := multipart.NewWriter(body)

	if path != "" {
		filePath := path
		file, err := os.Open(filePath)
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()

		w, err := mw.CreateFormFile(fieldName, filePath)
		if err != nil {
			t.Fatal(err)
		}

		if _, err := io.Copy(w, file); err != nil {
			t.Fatal(err)
		}
	}

	field, err := mw.CreateFormField("content")
	if err != nil {
		t.Fatal(err)
	}
	if content != "" {
		_, err = io.Copy(field, bytes.NewBufferString(content))
		if err != nil {
			t.Fatalf("Error writing to form field: %v", err)
		}
	}

	return body, mw
}

func TestIndexHandler(t *testing.T) {
	handler := setup()

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

func TestSearchPageHandler(t *testing.T) {
	handler := setup()

	t.Run("OK", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/search_page", nil)
		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("want: %v, got: %v", http.StatusOK, rec.Code)
		}
	})

	t.Run("err wrong method", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/search_page", nil)
		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("want: %v, got: %v", http.StatusMethodNotAllowed, rec.Code)
		}
	})
}

func TestSearchHandler(t *testing.T) {
	handler := setup()

	t.Run("OK", func(t *testing.T) {
		handler.Usecases.Posts.CreatePost(entity.Post{Title: "sdf"})
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/search", nil)

		form := url.Values{}
		form.Add("search", "abcd")
		req.PostForm = form

		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("want: %v, got: %v", http.StatusOK, rec.Code)
		}
	})

	t.Run("err wrong method", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/search", nil)
		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("want: %v, got: %v", http.StatusMethodNotAllowed, rec.Code)
		}
	})

	t.Run("err empty request", func(t *testing.T) {
		rec := httptest.NewRecorder()

		req := httptest.NewRequest(http.MethodPost, "/search", nil)

		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("want: %v, got: %v", http.StatusBadRequest, rec.Code)
		}
	})
}

func TestCreateCategoryPageHandler(t *testing.T) {
	handler := setup()

	t.Run("OK", func(t *testing.T) {
		if err := handler.Usecases.Users.SignUp(entity.User{}); err != nil {
			t.Fatal(err)
		}
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/create_category_page", nil)
		cookie := &http.Cookie{
			Name: "session_token",
		}
		req.AddCookie(cookie)

		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("want: %v, got: %v", http.StatusOK, rec.Code)
		}
	})

	t.Run("err not authorized", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/create_category_page", nil)
		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("want: %v, got: %v", http.StatusUnauthorized, rec.Code)
		}
	})

	t.Run("err method not allowed", func(t *testing.T) {
		if err := handler.Usecases.Users.SignUp(entity.User{}); err != nil {
			t.Fatal(err)
		}
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/create_category_page", nil)
		cookie := &http.Cookie{
			Name: "session_token",
		}
		req.AddCookie(cookie)

		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("want: %v, got: %v", http.StatusMethodNotAllowed, rec.Code)
		}
	})

	t.Run("err low access level", func(t *testing.T) {
		// if there is one user, he becomes admin and can create categories
		// if more, he behaves as simple user
		if err := handler.Usecases.Users.SignUp(entity.User{}); err != nil {
			t.Fatal(err)
		}

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/create_category_page", nil)
		cookie := &http.Cookie{
			Name: "session_token",
		}
		req.AddCookie(cookie)

		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Fatalf("want: %v, got: %v", http.StatusForbidden, rec.Code)
		}
	})
}

func TestCreateCategoryHandler(t *testing.T) {
	handler := setup()

	t.Run("OK", func(t *testing.T) {
		if err := handler.Usecases.Users.SignUp(entity.User{}); err != nil {
			t.Fatal(err)
		}
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/create_category", nil)

		cookie := &http.Cookie{
			Name: "session_token",
		}
		req.AddCookie(cookie)

		form := url.Values{}
		form.Add("category", "movies")
		req.PostForm = form

		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusFound {
			t.Fatalf("want: %v, got: %v", http.StatusFound, rec.Code)
		}
	})

	t.Run("err wrong method", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodHead, "/create_category", nil)

		cookie := &http.Cookie{
			Name: "session_token",
		}
		req.AddCookie(cookie)

		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("want: %v, got: %v", http.StatusMethodNotAllowed, rec.Code)
		}
	})

	t.Run("err empty request", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/create_category", nil)

		cookie := &http.Cookie{
			Name: "session_token",
		}
		req.AddCookie(cookie)

		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("want: %v, got: %v", http.StatusBadRequest, rec.Code)
		}
	})
}

func TestSearchByCategoryHandler(t *testing.T) {
	handler := setup()

	t.Run("OK", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/categories/cars", nil)
		if err := handler.Usecases.Posts.CreateCategories([]string{"cars"}); err != nil {
			t.Fatal(err)
		}
		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("want: %v, got: %v", http.StatusOK, rec.Code)
		}
	})

	t.Run("err wrong method", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/categories/", nil)

		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("want: %v, got: %v", http.StatusMethodNotAllowed, rec.Code)
		}
	})

	t.Run("err category not found", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/categories/qwerty", nil)
		if err := handler.Usecases.Posts.CreateCategories([]string{"cars"}); err != nil {
			t.Fatal(err)
		}
		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("want: %v, got: %v", http.StatusBadRequest, rec.Code)
		}
	})
}
