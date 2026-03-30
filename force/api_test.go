package force

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/oauth2"
)

func TestRefreshToken(t *testing.T) {
	tests := []struct {
		name string

		token *oauth2.Token

		error string

		code int
		body []byte
	}{
		{
			name: "successful request",

			token: &oauth2.Token{TokenType: "Bearer", AccessToken: "foo"},

			code: http.StatusOK,
			body: []byte(`{}`), // NOTE: mock "API resources" response
		},
		{
			name: "invalid session",

			token: &oauth2.Token{TokenType: "Bearer", AccessToken: "foo"},

			error: "failed to ensure token validity",

			code: http.StatusForbidden,
			body: []byte(`[{"errorCode": "INVALID_SESSION_ID"}]`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var initialized bool

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if initialized && r.URL.Path == "/services/data/v53.0" {
					// NOTE: mock error response for "API resources" request.
					w.WriteHeader(tt.code)
					_, _ = w.Write(tt.body)
					return
				}

				// NOTE: mock "successful" (albeit not "valid") responses for the
				// "API resources" and "sobjects" requests during initialization.

				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{}`))

				initialized = true
			})

			srv := httptest.NewServer(handler)
			t.Cleanup(srv.Close)

			source := oauth2.StaticTokenSource(tt.token)

			force, err := CreateWithTokenSource(source, srv.URL, DefaultAPIVersion, http.DefaultClient)

			if err != nil {
				t.Fatalf("failed to initialize client: err=%v", err)
			}

			refreshError := force.RefreshToken()

			switch {
			case tt.error == "" && refreshError != nil:
				t.Errorf("should be able to refresh token: err=%v", refreshError)
			case tt.error != "" && refreshError == nil:
				t.Errorf("should not be able to refresh token: want=%q", tt.error)
			case tt.error != "" && refreshError != nil && !strings.HasPrefix(refreshError.Error(), tt.error):
				t.Errorf("should return error with expected prefix: err=%v, want=%q", refreshError, tt.error)
			}
		})
	}
}
