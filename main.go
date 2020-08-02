package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.RequestID())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "TIME=${time_rfc3339} ID=${id} METHOD=${method} URI=${uri} STATUS=${status} LATENCY=${latency_human} BYTES_IN=${bytes_in} BYTES_OUT=${bytes_out}\n",
	}))

	authRoute := e.Group("/auth")
	authRoute.POST("/register", register)
	authRoute.POST("/login", login)

	clientRoute := e.Group("/clients")
	clientRoute.POST("/register", register)
	clientRoute.GET("/:id", clientDetail)
	clientRoute.POST("/user/register", register)
	clientRoute.POST("/user/login", login)

	managementRoute := e.Group("/management")
	managementRoute.GET("/clients", listClients)
	managementRoute.POST("/clients/revoke", revokeClient)

	go func() {
		if err := e.Start(":8888"); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func register(c echo.Context) error {

	type UserRegister struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var u UserRegister

	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	return c.JSON(http.StatusOK, struct {
		Message string `json:"message"`
	}{"Register success"})
}

func login(c echo.Context) error {
	return nil
}

func registerClient(c echo.Context) error {
	return nil
}

func clientDetail(c echo.Context) error {
	return nil
}

func listClients(c echo.Context) error {
	return nil
}

func revokeClient(c echo.Context) error {
	return nil
}
