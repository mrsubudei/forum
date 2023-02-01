package v1_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"forum/internal/entity"
)

func TestCreateCommentPageHandler(t *testing.T) {
	handler := setup()
	if err := handler.Usecases.Users.SignUp(entity.User{}); err != nil {
		t.Fatal(err)
	}

	t.Run("OK", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/create_comment_page/1", nil)
		cookie := &http.Cookie{
			Name: "session_token",
		}
		req.AddCookie(cookie)

		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("want: %v, got: %v", http.StatusOK, rec.Code)
		}
	})

	t.Run("err wrong method", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/create_comment_page/1", nil)
		cookie := &http.Cookie{
			Name: "session_token",
		}
		req.AddCookie(cookie)

		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("want: %v, got: %v", http.StatusMethodNotAllowed, rec.Code)
		}
	})

	t.Run("err wrong path", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/create_comment_page/dfg", nil)
		cookie := &http.Cookie{
			Name: "session_token",
		}
		req.AddCookie(cookie)

		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("want: %v, got: %v", http.StatusNotFound, rec.Code)
		}
	})
}

func TestCreateCommentHandler(t *testing.T) {
	handler := setup()
	if err := handler.Usecases.Users.SignUp(entity.User{}); err != nil {
		t.Fatal(err)
	}

	t.Run("OK", func(t *testing.T) {
		body, mw := CreateMultipartForm(t, "../../../../templates/img/github_auth_icon.jpg",
			"Lorem ipsum dolor sit amet.", "image")

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/create_comment/2", body)
		req.Header.Set("Content-Type", mw.FormDataContentType())
		cookie := &http.Cookie{
			Name: "session_token",
		}
		req.AddCookie(cookie)
		mw.Close()
		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusFound {
			t.Fatalf("want: %v, got: %v", http.StatusFound, rec.Code)
		}
	})

	t.Run("err wrong method", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/create_comment/2", nil)
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
		body, mw := CreateMultipartForm(t, "../../../../templates/img/buttons/home.png",
			"", "image")

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/create_comment/2", body)
		req.Header.Set("Content-Type", mw.FormDataContentType())
		cookie := &http.Cookie{
			Name: "session_token",
		}
		req.AddCookie(cookie)
		mw.Close()
		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("want: %v, got: %v", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("err wrong path", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/create_comment/dfg", nil)
		cookie := &http.Cookie{
			Name: "session_token",
		}
		req.AddCookie(cookie)

		form := url.Values{}
		form.Add("content", "Lorem ipsum dolor sit amet.")
		req.PostForm = form

		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("want: %v, got: %v", http.StatusNotFound, rec.Code)
		}
	})
}

func TestCommentPutLikeHandler(t *testing.T) {
	handler := setup()
	if err := handler.Usecases.Users.SignUp(entity.User{}); err != nil {
		t.Fatal(err)
	}

	t.Run("OK", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/put_comment_like/2", nil)
		cookie := &http.Cookie{
			Name: "session_token",
		}
		req.AddCookie(cookie)

		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusFound {
			t.Fatalf("want: %v, got: %v", http.StatusFound, rec.Code)
		}
	})

	t.Run("err wrong path", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/put_comment_like/wert", nil)
		cookie := &http.Cookie{
			Name: "session_token",
		}
		req.AddCookie(cookie)

		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("want: %v, got: %v", http.StatusNotFound, rec.Code)
		}
	})
}

func TestCommentPutDislikeHandler(t *testing.T) {
	handler := setup()
	if err := handler.Usecases.Users.SignUp(entity.User{}); err != nil {
		t.Fatal(err)
	}

	t.Run("OK", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/put_comment_dislike/2", nil)
		cookie := &http.Cookie{
			Name: "session_token",
		}
		req.AddCookie(cookie)

		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusFound {
			t.Fatalf("want: %v, got: %v", http.StatusFound, rec.Code)
		}
	})

	t.Run("err wrong path", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/put_comment_dislike/wert", nil)
		cookie := &http.Cookie{
			Name: "session_token",
		}
		req.AddCookie(cookie)

		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("want: %v, got: %v", http.StatusNotFound, rec.Code)
		}
	})
}
