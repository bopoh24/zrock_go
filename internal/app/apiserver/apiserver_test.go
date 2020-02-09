package apiserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bopoh24/zrock_go/internal/app/model"
	"github.com/bopoh24/zrock_go/internal/app/store/teststore"
	"github.com/stretchr/testify/assert"
)

func TestServer_handleHome(t *testing.T) {
	rec := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)

	srv := NewServer(teststore.New())
	srv.ServeHTTP(rec, req)
	assert.Equal(t, rec.Code, http.StatusOK)
	assert.Equal(t, rec.Body.String(), "Welcome to ZROCK API Server")
}

func TestServer_handleRegister(t *testing.T) {
	srv := NewServer(teststore.New())
	rec := httptest.NewRecorder()
	// check StatusMethodNotAllowed
	req, _ := http.NewRequest(http.MethodGet, "/api/auth/register", nil)
	srv.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)

	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"email":      "user@example.org",
				"password":   "Password1",
				"nickname":   "cool_hacker",
				"first_name": "Alexander",
				"last_name":  "Voronin",
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid payload",
			payload:      "",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid params",
			payload: map[string]string{
				"email": "invalid",
			},
			expectedCode: http.StatusUnprocessableEntity,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/api/auth/register", b)
			srv.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_handleLogin(t *testing.T) {

	srv := NewServer(teststore.New())
	rec := httptest.NewRecorder()
	// check StatusMethodNotAllowed
	req, _ := http.NewRequest(http.MethodGet, "/api/auth/login", nil)
	srv.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)

	badTestCases := []struct {
		name         string
		payload      interface{}
		errorText    string
		expectedCode int
	}{
		{
			name:         "no data",
			payload:      "",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "no username",
			payload: map[string]string{
				"password": "Password1",
			},
			errorText:    "username required",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "no password",
			payload: map[string]string{
				"username": "User123",
			},
			errorText:    "password required",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "invalid credentials",
			payload: map[string]string{
				"username": "User1",
				"password": "Password1",
			},
			errorText:    "incorrect username or password",
			expectedCode: http.StatusUnauthorized,
		},
	}

	respBody := make(map[string]string)

	for _, tc := range badTestCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/api/auth/login", b)
			srv.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
			if tc.errorText != "" {
				json.Unmarshal(rec.Body.Bytes(), &respBody)
				assert.Equal(t, tc.errorText, respBody["error"])
			}
		})
	}

	//Valid credentials
	user := model.UserFixture()
	user.Email = "name@mail.com"
	user.Password = "Password1"
	user.Nickname = "JohnDoe"
	srv.store.User().Create(user)

	goodTestCases := []struct {
		name    string
		payload map[string]string
	}{
		{
			name: "login by email",
			payload: map[string]string{
				"username": user.Email,
				"password": user.Password,
			},
		},
		{
			name: "login by nickname",
			payload: map[string]string{
				"username": user.Nickname,
				"password": user.Password,
			},
		},
	}

	for _, tc := range goodTestCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/api/auth/login", b)
			srv.ServeHTTP(rec, req)
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Contains(t, rec.Body.String(), "token")
		})
	}

}
