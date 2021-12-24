package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"example.com/model"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
  "golang.org/x/crypto/bcrypt"
)

var db *sql.DB

func main() {

	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "chatroom",
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected!")
	e := echo.New()
	e.GET("/users", getUser)
	e.POST("/users", register)
	e.Logger.Fatal(e.Start(":8080"))
}

func getUser(c echo.Context) error {
	return c.String(http.StatusOK, "Hello world")
}

func register(c echo.Context) error {
	u := new(model.User)
	if err := c.Bind(u); err != nil {
		return err
	}
  
  password, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 14)

	_, err := db.Exec("INSERT INTO users (name, email, password) values (?, ?, ?)", u.Name, u.Email, password)

	if err != nil {
		return err
	}
  
	return c.String(http.StatusOK, "Success")
}
