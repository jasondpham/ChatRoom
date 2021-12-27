package server

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
	"golang.org/x/crypto/bcrypt"
)

var secret string = "secret"

var Db *sql.DB

func ConnectToDb() {

	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "chatroom",
	}

	var err error
	Db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := Db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	fmt.Println("Connected!")
}

func GetUser(c echo.Context) error {
	user := make(map[string]string)

	if err := json.NewDecoder(c.Request().Body).Decode(&user); err != nil {
		return err
	}

	row := Db.QueryRow("SELECT * from users where email = ?", user["email"])

	entry := new(model.User)
	if err := row.Scan(&entry.Id, &entry.Name, &entry.Email, &entry.Password); err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword(entry.Password, []byte(user["password"])); err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "incorrect password",
		})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(entry.Id)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString([]byte(secret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Could not login",
		})
	}

	cookie := http.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
	}

	c.SetCookie(&cookie)
	return c.JSON(http.StatusOK, echo.Map{
		"message": "success",
	})
}

func Register(c echo.Context) error {
	u := make(map[string]string)
	if err := json.NewDecoder(c.Request().Body).Decode(&u); err != nil {
		return err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(u["password"]), 14)

	_, err := Db.Exec("INSERT INTO users (name, email, password) values (?, ?, ?)", u["name"], u["email"], password)

	if err != nil {
		return err
	}

	return c.String(http.StatusOK, "Success")
}

func USER(c echo.Context) error {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "unauthenticated",
		})
	}
	token, err := jwt.ParseWithClaims(cookie.Value, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "unauthenticated",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user model.User

	erro := Db.QueryRow("SELECT * FROM users WHERE id = ?", claims.Issuer).Scan(&user.Id, &user.Name, &user.Email, &user.Password)

	if erro != nil {
        return c.String(http.StatusBadRequest, "Not found")
    }

	return c.JSON(http.StatusOK, user)
}

func Logout(c echo.Context) error {
    cookie := http.Cookie {
        Name: "jwt",
        Value: "",
        Expires: time.Now().Add(-time.Hour),
        HttpOnly: true,
    }

    c.SetCookie(&cookie)
    return c.JSON(http.StatusOK, echo.Map {
        "message": "success",
    })
}
