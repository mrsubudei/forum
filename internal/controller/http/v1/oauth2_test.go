package v1_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOauthSigninHandler(t *testing.T) {
	handler := setup()
	t.Parallel()
	tests := []struct {
		method     string
		name       string
		url        string
		wantStatus int
	}{
		{
			name:       "OK google",
			wantStatus: http.StatusTemporaryRedirect,
			url:        "/oauth2_signin/google",
			method:     http.MethodGet,
		},
		{
			name:       "OK github",
			wantStatus: http.StatusTemporaryRedirect,
			url:        "/oauth2_signin/github",
			method:     http.MethodGet,
		},
		{
			name:       "OK mailru",
			wantStatus: http.StatusTemporaryRedirect,
			url:        "/oauth2_signin/mailru",
			method:     http.MethodGet,
		},
		{
			name:       "error wrong api",
			wantStatus: http.StatusNotFound,
			url:        "/oauth2_signin/alem",
			method:     http.MethodGet,
		},
		{
			name:       "error wrong method",
			wantStatus: http.StatusMethodNotAllowed,
			url:        "/oauth2_signin/google",
			method:     http.MethodPost,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(tt.method, tt.url, nil)
			handler.Mux.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Fatalf("want: %v, got: %v", tt.wantStatus, rec.Code)
			}
		})
	}
}
