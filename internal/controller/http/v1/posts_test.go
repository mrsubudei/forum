package v1_test

import (
	"forum/internal/entity"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestPostPageHandler(t *testing.T) {
	handler := setup()

	t.Run("OK", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/posts/2", nil)
		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("want: %v, got: %v", http.StatusOK, rec.Code)
		}
	})

	t.Run("err wrong method", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/posts/2", nil)
		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("want: %v, got: %v", http.StatusMethodNotAllowed, rec.Code)
		}
	})

	t.Run("err wrong path", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/posts/sefdfg", nil)
		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("want: %v, got: %v", http.StatusNotFound, rec.Code)
		}
	})
}

func TestCreatePostPageHandler(t *testing.T) {
	handler := setup()
	if err := handler.Usecases.Users.SignUp(entity.User{}); err != nil {
		t.Fatal(err)
	}

	t.Run("OK", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/create_post_page", nil)
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
		req := httptest.NewRequest(http.MethodPut, "/create_post_page", nil)
		cookie := &http.Cookie{
			Name: "session_token",
		}
		req.AddCookie(cookie)

		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("want: %v, got: %v", http.StatusMethodNotAllowed, rec.Code)
		}
	})
}

func TestCreatePostHandler(t *testing.T) {
	handler := setup()
	if err := handler.Usecases.Users.SignUp(entity.User{}); err != nil {
		t.Fatal(err)
	}

	t.Run("OK", func(t *testing.T) {
		body, mw := CreateMultipartForm(t, "",
			"Lorem ipsum dolor sit amet.", "")

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/create_post", body)
		req.Header.Set("Content-Type", mw.FormDataContentType())
		cookie := &http.Cookie{
			Name: "session_token",
		}
		req.AddCookie(cookie)
		form := url.Values{}
		form.Add("title", "BMW")
		form.Add("categories", "cars")
		req.PostForm = form
		mw.Close()
		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusFound {
			t.Fatalf("want: %v, got: %v", http.StatusFound, rec.Code)
		}
	})

	t.Run("err wrong method", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPatch, "/create_post", nil)
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
		body, mw := CreateMultipartForm(t, "",
			"", "")

		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/create_post", body)
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

	t.Run("err category is not chosen", func(t *testing.T) {
		body, mw := CreateMultipartForm(t, "",
			"abc", "")
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/create_post", body)
		req.Header.Set("Content-Type", mw.FormDataContentType())
		cookie := &http.Cookie{
			Name: "session_token",
		}
		req.AddCookie(cookie)
		mw.Close()
		form := url.Values{}
		form.Add("title", "BMW")
		form.Add("content", "Lorem ipsum dolor sit amet.")
		req.PostForm = form

		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("want: %v, got: %v", http.StatusBadRequest, rec.Code)
		}
	})
}

func TestPostPutLikeHandler(t *testing.T) {
	handler := setup()
	if err := handler.Usecases.Users.SignUp(entity.User{}); err != nil {
		t.Fatal(err)
	}

	t.Run("OK", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/put_post_like/1", nil)
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
		req := httptest.NewRequest(http.MethodPost, "/put_post_like/wewe", nil)
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

func TestPostPutDislikeHandler(t *testing.T) {
	handler := setup()
	if err := handler.Usecases.Users.SignUp(entity.User{}); err != nil {
		t.Fatal(err)
	}

	t.Run("OK", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/put_post_dislike/2", nil)
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
		req := httptest.NewRequest(http.MethodPost, "/put_post_dislike/12/e", nil)
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

func TestFindPostsHandler(t *testing.T) {
	handler := setup()
	if err := handler.Usecases.Users.SignUp(entity.User{}); err != nil {
		t.Fatal(err)
	}

	t.Run("OK", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/find_posts/liked/1", nil)
		cookie := &http.Cookie{
			Name: "session_token",
		}
		req.AddCookie(cookie)

		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("want: %v, got: %v", http.StatusOK, rec.Code)
		}
	})

	t.Run("err wrogn method", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodOptions, "/find_posts/1", nil)
		cookie := &http.Cookie{
			Name: "session_token",
		}
		req.AddCookie(cookie)

		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("want: %v, got: %v", http.StatusMethodNotAllowed, rec.Code)
		}
	})

	t.Run("err wrong query", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/find_posts/1", nil)
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
