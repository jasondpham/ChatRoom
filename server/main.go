package main

import (
	"database/sql"
	"encoding/json"
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
  user := make(map[string]string)

  if err := json.NewDecoder(c.Request().Body).Decode(&user); err != nil {
    return err
  }

  row := db.QueryRow("SELECT * from users where email = ?", user["email"])

  entry := new(model.User)
  if err := row.Scan(&entry.Id, &entry.Name, &entry.Email, &entry.Password); err != nil {
    return err
  }

  if err := bcrypt.CompareHashAndPassword(entry.Password, []byte(user["password"])); err != nil {
    return c.JSON(http.StatusUnauthorized, echo.Map {
      "message": "incorrect password",
    })
  }
  return c.JSON(http.StatusOK, echo.Map {
    "name": entry.Name,
  })
}

func register(c echo.Context) error {
	u := make(map[string]string) 
	if err := json.NewDecoder(c.Request().Body).Decode(&u); err != nil {
		return err
	}
  
  password, _ := bcrypt.GenerateFromPassword([]byte(u["password"]), 14)

	_, err := db.Exec("INSERT INTO users (name, email, password) values (?, ?, ?)", u["name"], u["email"], password)

	if err != nil {
		return err
	}
  
	return c.String(http.StatusOK, "Success")
}
