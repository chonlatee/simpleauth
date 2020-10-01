package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestRouteRegisterSuccess(t *testing.T) {
	e := echo.New()
	u := &mockUserRepo{
		insertFunc: func(ctx context.Context, username, email, password string) (int64, error) {
			return 0, nil
		},
	}

	userJSON := `{"username":"chonlatee","password": "123456", "email":"jon@labstack.com"}`
	req, _ := http.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(userJSON))
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

	if rec.Code != http.StatusOK {
		t.Errorf("Expect status code equal %d but got %d", http.StatusOK, rec.Code)
	}

	if rr.Message != "Register success" {
		t.Errorf("Expect response message equal `Register success` but got `%s`", rr.Message)
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
