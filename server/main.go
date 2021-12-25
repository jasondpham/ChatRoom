package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"example.com/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var secret string = "secret"

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

  e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
    AllowCredentials: true,
  }))
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

  claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims {
    Issuer: strconv.Itoa(int(entry.Id)),
    ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
  })

  token, err := claims.SignedString([]byte(secret))
  if err != nil {
    return c.JSON(http.StatusInternalServerError, echo.Map {
      "message": "Could not login",
    })
  }

  cookie := http.Cookie {
    Name: "jwt",
    Value: token,
    Expires: time.Now().Add(time.Hour * 24),
    HttpOnly: true,
  }

  c.SetCookie(&cookie)
  return c.JSON(http.StatusOK, echo.Map {
    "message": "success",
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
