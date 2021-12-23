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
	e.POST("/users", addUser)
	e.Logger.Fatal(e.Start(":8080"))
}

func getUser(c echo.Context) error {
	return c.String(http.StatusOK, "Hello world")
}

func addUser(c echo.Context) error {
	// u := make(map[string]string)
	// if err := json.NewDecoder(c.Request().Body).Decode(&u); err != nil {
	// 	return err
	// }

	u := new(model.User)
	if err := c.Bind(u); err != nil {
		return err
	}
	// user := model.User{
	// 	Name:     u["name"],
	// 	Email:    u["email"],
	// 	Password: u["password"],
	// }

	_, err := db.Exec("INSERT INTO users (name, email, password) values (?, ?, ?)", u.Name, u.Email, u.Password)

	if err != nil {
		return err
	}

	return c.String(http.StatusOK, "Success")
}
