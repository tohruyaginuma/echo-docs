package main

import (
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

type User struct {
	Name string `json:"name" xml:"name" form:"name" query:"name"`
	Email string `json:"email" xml:"email" form:"email" query:"email"`
}

func home (c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func show (c echo.Context) error {
	team := c.QueryParam("team")
	member := c.QueryParam("member")

	return c.String(http.StatusOK, "team:" + team + ", member:" + member)
}

func save (c echo.Context) error {
	name := c.FormValue("name")
	avatar, err := c.FormFile("avatar")

	if err != nil {
		return err
	}
	
	src, err := avatar.Open()

	if err != nil {
		return err
	}

	defer src.Close()

	dst, err := os.Create(avatar.Filename)

	if err != nil {
		return err
	}

	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return c.String(http.StatusOK, "<b>Thank you! " + name + "</b>")
}

func getUser (c echo.Context) error {
	id := c.Param("id")
	return c.String(http.StatusOK, id)
}

func saveUser (c echo.Context) error {
	user := new(User)
	if err := c.Bind(user); err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, user)
}

func updateUser (c echo.Context) error {
	return c.String(http.StatusOK, "TODO: Update User!")
}

func deleteUser (c echo.Context) error {
	return c.String(http.StatusOK, "TODO: DELETE User!")
}

func main () {
	e := echo.New()

	e.GET("/", home)
	e.GET("/show", show)

	e.POST("/save", save)

	e.POST("/users", saveUser)

	e.GET("/users/:id", getUser)
	e.PUT("/users/:id", updateUser)
	e.DELETE("/users/:id", deleteUser)

	e.Logger.Fatal(e.Start(":1323"))
}