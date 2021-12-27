package main

import (
	"database/sql"

	"example.com/server"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var db *sql.DB
var secret string = "secret"

func main() {

  server.ConnectToDb()
  db = server.Db
	e := echo.New()

  e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowCredentials: true,
  }))
	e.GET("/login", server.GetUser)
	e.POST("/users", server.Register)
    e.GET("/users", server.USER)
	e.Logger.Fatal(e.Start(":8080"))
}

