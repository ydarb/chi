package middleware

import (
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCleanPath(t *testing.T) {
	r := chi.NewRouter()
	r.Use(CleanPath)
	r.Get("/test/path", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(chi.RouteContext(r.Context()).RoutePath))
	})
	r.Connect("/*", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(chi.RouteContext(r.Context()).RoutePath))
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	tests := []struct {
		name         string
		method       string
		path         string
		expectedPath string
		expectedCode int
	}{
		{
			name:         "get request clean path",
			method:       http.MethodGet,
			path:         "/test/path///",
			expectedPath: "/test/path",
			expectedCode: http.StatusOK,
		},
		{
			name:         "connect request do nothing with path",
			method:       http.MethodConnect,
			path:         "/test/path///",
			expectedPath: "",
			expectedCode: http.StatusOK,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			resp, got := testRequestWithCleanPath(t, ts, tc.method, tc.path)
			if resp.StatusCode != tc.expectedCode {
				t.Errorf("expected status %d, got %d", tc.expectedCode, resp.StatusCode)
			}
			if got != tc.expectedPath {
				t.Errorf("expected path %q but got %q", tc.expectedPath, got)
			}
		})
	}
}

func testRequestWithCleanPath(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return nil, ""
	}
	defer resp.Body.Close()
	return resp, string(respBody)
}
