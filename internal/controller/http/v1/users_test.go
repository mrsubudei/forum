package v1_test

import (
	"forum/internal/entity"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestUserPageHandler(t *testing.T) {
	handler := setup()

	t.Run("OK", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/users/3", nil)
		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("want: %v, got: %v", http.StatusOK, rec.Code)
		}
	})

	t.Run("err wrong method", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/users/3", nil)
		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("want: %v, got: %v", http.StatusMethodNotAllowed, rec.Code)
		}
	})

	t.Run("err wrong path", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/users/3/2", nil)
		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("want: %v, got: %v", http.StatusNotFound, rec.Code)
		}
	})
}

func TestAllUsersPageHandler(t *testing.T) {
	handler := setup()

	t.Run("OK", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/all_users_page", nil)
		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("want: %v, got: %v", http.StatusOK, rec.Code)
		}
	})

	t.Run("err wrong method", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodHead, "/all_users_page", nil)
		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("want: %v, got: %v", http.StatusMethodNotAllowed, rec.Code)
		}
	})
}

func TestSignUpPageHandler(t *testing.T) {
	handler := setup()

	t.Run("OK", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/signup_page", nil)
		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("want: %v, got: %v", http.StatusOK, rec.Code)
		}
	})

	t.Run("err wrong method", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodOptions, "/signup_page", nil)
		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("want: %v, got: %v", http.StatusMethodNotAllowed, rec.Code)
		}
	})
}

func TestSignUpHandler(t *testing.T) {
	handler := setup()

	t.Run("OK", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/signup", nil)

		form := url.Values{}
		form.Add("user", "Riddle")
		form.Add("email", "Riddle@mail.ru")
		form.Add("password", "Vivse")
		form.Add("confirm_password", "Vivse")
		req.PostForm = form

		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusFound {
			t.Fatalf("want: %v, got: %v", http.StatusFound, rec.Code)
		}
	})

	t.Run("err wrong method", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPut, "/signup", nil)

		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("want: %v, got: %v", http.StatusMethodNotAllowed, rec.Code)
		}
	})

	t.Run("err empty name", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/signup", nil)

		form := url.Values{}
		form.Add("email", "Riddle@mail.ru")
		form.Add("password", "Vivse")
		form.Add("confirm_password", "Vivse")
		req.PostForm = form

		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("want: %v, got: %v", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("err empty email", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/signup", nil)

		form := url.Values{}
		form.Add("user", "Riddle")
		form.Add("password", "Vivse")
		form.Add("confirm_password", "Vivse")
		req.PostForm = form

		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("want: %v, got: %v", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("err empty password", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/signup", nil)

		form := url.Values{}
		form.Add("user", "Riddle")
		form.Add("email", "Riddle@mail.ru")
		form.Add("confirm_password", "Vivse")
		req.PostForm = form

		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("want: %v, got: %v", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("err empty confirm password", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/signup", nil)

		form := url.Values{}
		form.Add("user", "Riddle")
		form.Add("email", "Riddle@mail.ru")
		form.Add("password", "Vivse")
		req.PostForm = form

		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("want: %v, got: %v", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("err incorrect email", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/signup", nil)

		form := url.Values{}
		form.Add("user", "Riddle")
		form.Add("email", "Riddletrt")
		form.Add("password", "Vivse")
		form.Add("confirm_password", "Vivse")
		req.PostForm = form

		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("want: %v, got: %v", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("err passwords are not same", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/signup", nil)

		form := url.Values{}
		form.Add("user", "Riddle")
		form.Add("email", "Riddle@mail.ru")
		form.Add("password", "Vivse")
		form.Add("confirm_password", "notsame")
		req.PostForm = form

		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("want: %v, got: %v", http.StatusBadRequest, rec.Code)
		}
	})

}

func TestSignInPageHandler(t *testing.T) {
	handler := setup()

	t.Run("OK", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/signin_page", nil)
		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("want: %v, got: %v", http.StatusOK, rec.Code)
		}
	})

	t.Run("err wrong method", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodTrace, "/signin_page", nil)
		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("want: %v, got: %v", http.StatusMethodNotAllowed, rec.Code)
		}
	})
}

func TestSignInHandler(t *testing.T) {
	handler := setup()

	t.Run("OK", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/signin", nil)

		form := url.Values{}
		form.Add("user", "Riddle")
		form.Add("password", "Vivse")
		req.PostForm = form

		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusFound {
			t.Fatalf("want: %v, got: %v", http.StatusFound, rec.Code)
		}
	})

	t.Run("err wrong method", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodConnect, "/signin", nil)

		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("want: %v, got: %v", http.StatusMethodNotAllowed, rec.Code)
		}
	})

	t.Run("err empty name", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/signin", nil)

		form := url.Values{}
		form.Add("password", "Vivse")
		req.PostForm = form

		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("want: %v, got: %v", http.StatusBadRequest, rec.Code)
		}
	})

	t.Run("err empty password", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/signin", nil)

		form := url.Values{}
		form.Add("user", "Riddle")
		req.PostForm = form

		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("want: %v, got: %v", http.StatusBadRequest, rec.Code)
		}
	})
}

func TestEditProfilePageHandler(t *testing.T) {
	handler := setup()
	if err := handler.Usecases.Users.SignUp(entity.User{}); err != nil {
		t.Fatal(err)
	}

	t.Run("OK", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/edit_profile_page/1", nil)
		cookie := &http.Cookie{
			Name: "session_token",
		}
		req.AddCookie(cookie)

		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("want: %v, got: %v", http.StatusOK, rec.Code)
		}
	})

	t.Run("wrong method", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodDelete, "/edit_profile_page/1", nil)
		cookie := &http.Cookie{
			Name: "session_token",
		}
		req.AddCookie(cookie)

		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Fatalf("want: %v, got: %v", http.StatusMethodNotAllowed, rec.Code)
		}
	})

	t.Run("err wrong pass", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/edit_profile_page/dfg", nil)
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

func TestEditProfileHandler(t *testing.T) {
	handler := setup()
	if err := handler.Usecases.Users.SignUp(entity.User{}); err != nil {
		t.Fatal(err)
	}

	t.Run("OK", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/edit_profile/", nil)
		cookie := &http.Cookie{
			Name: "session_token",
		}
		req.AddCookie(cookie)

		form := url.Values{}
		form.Add("id", "3")
		form.Add("city", "Astana")
		form.Add("sign", "new sign")
		req.PostForm = form

		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusFound {
			t.Fatalf("want: %v, got: %v", http.StatusFound, rec.Code)
		}
	})

	t.Run("wrong method", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/edit_profile/", nil)
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

func TestSignOutHandler(t *testing.T) {
	handler := setup()
	if err := handler.Usecases.Users.SignUp(entity.User{}); err != nil {
		t.Fatal(err)
	}

	t.Run("OK", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/signout", nil)
		cookie := &http.Cookie{
			Name: "session_token",
		}
		req.AddCookie(cookie)

		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusFound {
			t.Fatalf("want: %v, got: %v", http.StatusFound, rec.Code)
		}
	})
}

func TestFindReactedUsersHandler(t *testing.T) {
	handler := setup()
	if err := handler.Usecases.Users.SignUp(entity.User{}); err != nil {
		t.Fatal(err)
	}

	t.Run("OK liked post", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/find_reacted_users/post/liked/4", nil)
		cookie := &http.Cookie{
			Name: "session_token",
		}
		req.AddCookie(cookie)

		handler.Mux.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("want: %v, got: %v", http.StatusOK, rec.Code)
		}
	})

	t.Run("OK disliked comment", func(t *testing.T) {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/find_reacted_users/comment/disliked/3", nil)
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
		req := httptest.NewRequest(http.MethodPut, "/find_reacted_users/comment/disliked/3", nil)
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
		req := httptest.NewRequest(http.MethodGet, "/find_reacted_users/comment/3", nil)
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
