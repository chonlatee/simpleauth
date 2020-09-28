package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/chonlatee/authserver/pkg/models/mariadb"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	e := echo.New()

	e.Use(middleware.RequestID())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "TIME=${time_rfc3339} ID=${id} METHOD=${method} URI=${uri} STATUS=${status} LATENCY=${latency_human} BYTES_IN=${bytes_in} BYTES_OUT=${bytes_out}\n",
	}))

	go func() {
		db := openDB()

		r := route{
			userModel: &mariadb.UserModel{DB: db},
		}

		authRoute := e.Group("/auth")
		authRoute.POST("/register", r.register)
		authRoute.POST("/login", r.login)

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
