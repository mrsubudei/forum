package v1_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	v1 "forum/internal/controller/http/v1"
	"forum/internal/entity"
)

func getMockHandlerOne(t *testing.T, key string) http.HandlerFunc {
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		content, ok := r.Context().Value(v1.Key(key)).(v1.Content)
		if !ok {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
		if !content.Authorized {
			t.Fatalf("want: true, got: false")
		}
		if content.User.Id == 1 && !content.Admin {
			t.Fatalf("want: true, got: false")
		}
		if !content.Admin {
			w.WriteHeader(http.StatusForbidden)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	})
	return mockHandler
}

func getMockHandlerTwo(t *testing.T, key string) http.HandlerFunc {
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		content, ok := r.Context().Value(v1.Key(key)).(v1.Content)
		if !ok {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
		if content.Authorized && content.Unauthorized || !content.Authorized && !content.Unauthorized {
			t.Fatalf("should not both be true/false")
		}
		if content.User.Id == 1 && !content.Admin {
			t.Fatalf("want: true, got: false")
		}
		if content.Admin && !content.Authorized {
			t.Fatalf("want: true, got: false")
		}
		w.WriteHeader(http.StatusOK)
	})
	return mockHandler
}

func TestCheckAuth(t *testing.T) {
	handler := setup()

	t.Run("OK", func(t *testing.T) {
		if err := handler.Usecases.Users.SignUp(entity.User{}); err != nil {
			t.Fatal(err)
		}

		mockHandler := getMockHandlerOne(t, "content")

		handlerToTest := handler.CheckAuth(mockHandler)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
		cookie := &http.Cookie{
			Name: "session_token",
		}

		req.AddCookie(cookie)
		handlerToTest.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("want: %v, got: %v", http.StatusOK, rec.Code)
		}
	})

	t.Run("err user is not admin (low access level)", func(t *testing.T) {
		// creating another user and id will become not 1
		if err := handler.Usecases.Users.SignUp(entity.User{}); err != nil {
			t.Fatal(err)
		}

		mockHandler := getMockHandlerOne(t, "content")

		handlerToTest := handler.CheckAuth(mockHandler)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
		cookie := &http.Cookie{
			Name: "session_token",
		}

		req.AddCookie(cookie)
		handlerToTest.ServeHTTP(rec, req)

		if rec.Code != http.StatusForbidden {
			t.Fatalf("want: %v, got: %v", http.StatusForbidden, rec.Code)
		}
	})

	t.Run("err not authorized", func(t *testing.T) {
		mockHandler := getMockHandlerOne(t, "content")

		handlerToTest := handler.CheckAuth(mockHandler)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "http://testing", nil)

		handlerToTest.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("want: %v, got: %v", http.StatusUnauthorized, rec.Code)
		}
	})

	t.Run("err wrong context", func(t *testing.T) {
		mockHandler := getMockHandlerOne(t, "wrongKey")

		handlerToTest := handler.CheckAuth(mockHandler)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
		cookie := &http.Cookie{
			Name: "session_token",
		}

		req.AddCookie(cookie)

		handlerToTest.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotAcceptable {
			t.Fatalf("want: %v, got: %v", http.StatusNotAcceptable, rec.Code)
		}
	})
}

func TestAssignStatus(t *testing.T) {
	handler := setup()

	t.Run("OK", func(t *testing.T) {
		if err := handler.Usecases.Users.SignUp(entity.User{}); err != nil {
			t.Fatal(err)
		}

		mockHandler := getMockHandlerTwo(t, "content")

		handlerToTest := handler.AssignStatus(mockHandler)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
		cookie := &http.Cookie{
			Name: "session_token",
		}

		req.AddCookie(cookie)
		handlerToTest.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("want: %v, got: %v", http.StatusOK, rec.Code)
		}
	})

	t.Run("OK not auth", func(t *testing.T) {
		mockHandler := getMockHandlerTwo(t, "content")

		handlerToTest := handler.AssignStatus(mockHandler)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "http://testing", nil)

		handlerToTest.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("want: %v, got: %v", http.StatusOK, rec.Code)
		}
	})

	t.Run("err wrong context", func(t *testing.T) {
		mockHandler := getMockHandlerTwo(t, "wrongKey")

		handlerToTest := handler.AssignStatus(mockHandler)

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "http://testing", nil)
		cookie := &http.Cookie{
			Name: "session_token",
		}

		req.AddCookie(cookie)

		handlerToTest.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotAcceptable {
			t.Fatalf("want: %v, got: %v", http.StatusNotAcceptable, rec.Code)
		}
	})
}
