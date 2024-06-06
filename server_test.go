package yadro

import (
	"bytes"
	"encoding/json"
	"github.com/golang-jwt/jwt/v4"
	_ "github.com/golang-jwt/jwt/v4"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	http2 "yadro/internal/adapter/handler/http"
)

var jwtKey = []byte("my_secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func generateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		name               string
		username           string
		expectedStatusCode int
	}{
		{
			name:               "Successful update",
			username:           "admin",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Non-admin update",
			username:           "user",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "Unauthorised user",
			username:           "",
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/update", nil)

			if tt.username != "" {
				token, err := generateJWT(tt.username)
				if err != nil {
					t.Fatalf("Failed to generate a current: %v", err)
				}
				cookie := &http.Cookie{
					Name:  "token",
					Value: token,
				}
				req.AddCookie(cookie)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(http2.Update)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatusCode)
			}
		})
	}
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

func TestLogin(t *testing.T) {
	tests := []struct {
		name               string
		username           string
		password           string
		expectedStatusCode int
	}{
		{
			name:               "Successful login",
			username:           "admin",
			password:           "admin",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Invalid username",
			username:           "wronguser",
			password:           "admin",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "Invalid password",
			username:           "admin",
			password:           "wrongpassword",
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			creds := Credentials{
				Username: tt.username,
				Password: tt.password,
			}
			body, _ := json.Marshal(creds)
			req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(http2.Login)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatusCode)
			}

			if tt.expectedStatusCode == http.StatusOK {
				cookie := rr.Result().Cookies()
				if len(cookie) == 0 || cookie[0].Name != "token" {
					t.Errorf("Expected a token cookie, but none was found")
				}
			}
		})
	}
}

func TestPics(t *testing.T) {
	tests := []struct {
		name               string
		username           string
		tokenValid         bool
		expectedStatusCode int
	}{
		{
			name:               "Invalid token",
			username:           "user",
			tokenValid:         false,
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			name:               "Ok",
			username:           "admin",
			tokenValid:         true,
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenString := ""
			if tt.tokenValid {
				var err error
				tokenString, err = generateJWT(tt.username)
				if err != nil {
					t.Fatalf("Не удалось создать токен: %v", err)
				}
			} else {
				tokenString = "invalid_token"
			}

			req, err := http.NewRequest(http.MethodGet, "/pics", nil)
			if err != nil {
				t.Fatalf("Не удалось создать запрос: %v", err)
			}
			req.AddCookie(&http.Cookie{
				Name:  "token",
				Value: tokenString,
			})
			q := req.URL.Query()
			q.Add("search", "example")
			req.URL.RawQuery = q.Encode()

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(http2.Pics)
			handler.ServeHTTP(rr, req)
			if status := rr.Code; status != tt.expectedStatusCode {
				t.Errorf("Неправильный статус код: получил %v, ожидал %v", status, tt.expectedStatusCode)
			}

		})
	}
}
