package identity

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractor(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v, ok := r.Context().Value(IDHeaderKey).(*Identity)
		require.True(t, ok)
		require.Equal(t, "user", v.User)
	})

	resp := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://test", nil)
	Extractor(testHandler).ServeHTTP(resp, req)
	require.Equal(t, 400, resp.Code)

	resp = httptest.NewRecorder()
	req.Header[FedoraIDHeader] = []string{"garbage"}
	Extractor(testHandler).ServeHTTP(resp, req)
	require.Equal(t, 400, resp.Code)

	resp = httptest.NewRecorder()
	req.Header[FedoraIDHeader] = []string{"	eyJ1c2VyIjoiIn0K"}
	Extractor(testHandler).ServeHTTP(resp, req)
	require.Equal(t, 400, resp.Code)

	resp = httptest.NewRecorder()
	req.Header[FedoraIDHeader] = []string{"eyJ1c2VyIjoidXNlciJ9"}
	Extractor(testHandler).ServeHTTP(resp, req)
	require.Equal(t, 200, resp.Code)
}
