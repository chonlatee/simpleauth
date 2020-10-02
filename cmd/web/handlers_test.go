package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func genChar(length int) string {
	pass := ""
	for i := 0; i < length; i++ {
		pass += "a"
	}
	return pass
}

func TestValidateRegisterPayload(t *testing.T) {

	shortString := genChar(5)
	longString := genChar(51)

	var registerTestCase = []struct {
		name          string
		path          string
		payload       string
		resultMessage string
		statusCode    int
	}{
		{
			name:          "Register success",
			path:          "/auth/register",
			payload:       `{"username":"chonlatee","password": "123456", "email":"jon@labstack.com"}`,
			resultMessage: "Register success",
			statusCode:    http.StatusOK,
		},
		{
			name:          "Register fail with password length is too short",
			path:          "/auth/register",
			payload:       fmt.Sprintf(`{"username":"chonlatee","password": "%s", "email":"jon@labstack.com"}`, shortString),
			resultMessage: "Invalid password length",
			statusCode:    http.StatusBadRequest,
		},
		{
			name:          "Register fail with password length is too long",
			path:          "/auth/register",
			payload:       fmt.Sprintf(`{"username":"chonlatee","password": "%s", "email":"jon@labstack.com"}`, longString),
			resultMessage: "Invalid password length",
			statusCode:    http.StatusBadRequest,
		},
		{
			name:          "Register fail with username length is too short",
			path:          "/auth/register",
			payload:       fmt.Sprintf(`{"username":"%s","password": "123456", "email":"jon@labstack.com"}`, shortString),
			resultMessage: "Invalid username length",
			statusCode:    http.StatusBadRequest,
		},
		{
			name:          "Register fail with username length is too long",
			path:          "/auth/register",
			payload:       fmt.Sprintf(`{"username":"%s","password": "123456", "email":"jon@labstack.com"}`, longString),
			resultMessage: "Invalid username length",
			statusCode:    http.StatusBadRequest,
		},
		{
			name:          "Register fail with invalid email",
			path:          "/auth/register",
			payload:       `{"username":"chonlatee","password": "123456", "email":"foo"}`,
			resultMessage: "Invalid email",
			statusCode:    http.StatusBadRequest,
		},
	}

	e := echo.New()
	u := &mockUserRepo{
		insertFunc: func(ctx context.Context, username, email, password string) (int64, error) {
			return 0, nil
		},
	}

	for _, tc := range registerTestCase {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPost, tc.path, strings.NewReader(tc.payload))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()

			c := e.NewContext(req, rec)

			r := NewRoute(u)

			err := r.register(c)

			if err != nil {
				t.Error("expect error equal nil")
			}

			type Result struct {
				Message string `json:"message"`
			}

			var rr Result

			json.Unmarshal(rec.Body.Bytes(), &rr)

			if rec.Code != tc.statusCode {
				t.Errorf("Expect status code equal %d but got %d", tc.statusCode, rec.Code)
			}

			if rr.Message != tc.resultMessage {
				t.Errorf("Expect response message equal `%s` but got `%s`", tc.resultMessage, rr.Message)
			}
		})
	}

}

type mockUserRepo struct {
	insertFunc     func(context.Context, string, string, string) (int64, error)
	getByEmailFunc func(context.Context, string) (string, error)
}

func (m *mockUserRepo) Insert(ctx context.Context, username, email, password string) (int64, error) {
	return m.insertFunc(ctx, username, email, password)
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (string, error) {
	return m.getByEmailFunc(ctx, email)
}
