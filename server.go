package main

import (
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/middleware"
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

func auth(username, password string, c echo.Context) (bool, error) {
	if username == "joe" && password == "secret" {
		return true, nil
	}

	return false, nil
}



func main () {
	e := echo.New()

	// Root Level middle ware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Group level middle ware
	g:= e.Group("/admin")
	g.Use(middleware.BasicAuth(auth))

	// Route level middle ware
	track := func (next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			println("request to /users")
			return next(c)
		}
	}


	e.GET("/", home, track)
	e.GET("/show", show, track)

	e.POST("/save", save, track)

	e.POST("/users", saveUser, track)

	e.GET("/users/:id", getUser, track)
	e.PUT("/users/:id", updateUser, track)
	e.DELETE("/users/:id", deleteUser, track)

	e.Static("/static", "static")

	e.Logger.Fatal(e.Start(":1323"))
}