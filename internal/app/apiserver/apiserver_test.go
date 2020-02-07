package apiserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bopoh24/zrock_go/internal/app/store/teststore"
	"github.com/stretchr/testify/assert"
)

func TestServer_handleHome(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)

	srv := NewServer(NewConfig(), teststore.New())
	srv.ServeHTTP(rec, req)
	assert.Equal(t, rec.Code, http.StatusOK)
	assert.Equal(t, rec.Body.String(), "Welcome to ZROCK API Server")
}
