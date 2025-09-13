package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type UserDTO struct {
	Name  string `json:"name" form:"name" query:"name"`
	Email string `json:"email" form:"email" query:"email"`
}

// type User struct {
// 	Name    string
// 	Email   string
// 	IsAdmin bool
// }

type CustomContext struct {
	echo.Context
}

type (
	User struct {
	  Name  string `json:"name" validate:"required"`
	  Email string `json:"email" validate:"required,email"`
	}
  
	CustomValidator struct {
	  validator *validator.Validate
	}
  )

  func (cv *CustomValidator) Validate(i interface{}) error {
	
	// バインドされた構造体のチェック
	if err := cv.validator.Struct(i); err != nil {
	  // Optionally, you could return the error to give each route more control over the status code
	  return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
  }

func (c *CustomContext) Foo() {
	println("foo")
}

func (c *CustomContext) Bar() {
	println("bar")
}

func home (c echo.Context) error {
	cc := c.(*CustomContext)
	cc.Foo()
	cc.Bar()
	return cc.String(200, "OK")
	// return c.String(http.StatusOK, "Hello, World!")
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

func getUsers(c echo.Context) (err error) {
	// u := new(UserDTO)
	// if err := c.Bind(u); err != nil {
	//   return c.String(http.StatusBadRequest, "bad request")
	// }

	c.Response().Before(func() {
		println("before response")
	})
	c.Response().After(func() {
	println("after response")
	})

	// Load into separate struct for security
	// user := User{
	//   Name: u.Name,
	//   Email: u.Email,
	//   IsAdmin: false,
	// }

	// executeSomeBusinessLogic(user)

	return c.String(http.StatusOK, "OK")
}

func fluentBinding(c echo.Context) (err error) {
	// url =  "/api/search?active=true&id=1&id=2&id=3&length=25"
	var opts struct {
		IDs []int64
		Active bool
	}
	length := int64(50) // default length is 50
	
	// creates query params binder that stops binding at first error
	err = echo.QueryParamsBinder(c).
		Int64("length", &length).
		Int64s("ids", &opts.IDs).
		Bool("active", &opts.Active).
		BindError() // returns first binding error

	return c.JSON(http.StatusOK, err)
}

func writeCookie(c echo.Context) error{
	cookie := new(http.Cookie)
	cookie.Name = "username"
	cookie.Value = "jon"
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)

	return c.String(http.StatusOK, "write a cookie")
}

func readCookie(c echo.Context) error {
	cookie, err := c.Cookie("username")
	if err != nil {
		return err
	}

	fmt.Println(cookie.Name)
	fmt.Println(cookie.Value)

	return c.String(http.StatusOK, "read a cookie")
}

func readAllCookies(c echo.Context) error {
	for _, cookie := range c.Cookies() {
		fmt.Println(cookie.Name)
		fmt.Println(cookie.Value)
	}

	return c.String(http.StatusOK, "read all the cookies")
}

func customHTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	code := http.StatusInternalServerError
	if he, ok :=  err.(*echo.HTTPError); ok {
		code = he.Code
	}

	c.Logger().Error(err)

	errorPage := fmt.Sprintf("%d.html", code)

	if err := c.File(errorPage); err != nil {
		c.Logger().Error(err)
	}

	
}


func main () {
	e := echo.New()

	e.Debug = true
	// e.HideBanner = true
	// e.HidePort = true
	// e.Pre(middleware.RemoveTrailingSlash())
	// Custom http error handler
	// e.HTTPErrorHandler = customHTTPErrorHandler

	// カスタムバリデータ登録できる
	// e.Validator = &CustomValidator{validator: validator.New()}

	// e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
	// 	return func(c echo.Context) error {
	// 		cc := &CustomContext{c}
	// 		return next(cc)
	// 	}
	// })

	// e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
	// 	return func(c echo.Context) error {
	// 		return echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid credentials")
	// 	}
	// })

	// Root Level middle ware
	// e.Use(middleware.Logger())
	// e.Use(middleware.Recover())

	// Group level middle ware
	// g:= e.Group("/admin")
	// g.Use(middleware.BasicAuth(auth))

	// Route level middle ware
	track := func (next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// println("request to /users")
			return next(c)
		}
	}


	e.GET("/", home, track)
	e.GET("/show", show, track)
	
	e.GET("/fluent", fluentBinding, track)

	e.POST("/save", save, track)

	e.GET("/users", getUsers, track)
	// e.POST("/users", saveUser, track)
	e.POST("/users", func(c echo.Context) (err error) {
		u := new(User)
		if err = c.Bind(u); err != nil {
		  return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err = c.Validate(u); err != nil {
		  return err
		}
		return c.JSON(http.StatusOK, u)
	  })

	e.GET("/users/:id", getUser, track)
	e.PUT("/users/:id", updateUser, track)
	e.DELETE("/users/:id", deleteUser, track)

	e.Static("/static", "static")

	// Echo logger
	e.Logger.Fatal(e.Start(":1323"))


	// if err := e.Start(":8080"); err != http.ErrServerClosed {
	// 	log.Fatal(err)
	//   }

	// s := http.Server{
	// 	Addr: ":8080",
	// 	Handler: e,
	// }

	// if err := s.ListenAndServe(); err != http.ErrServerClosed {
	// 	// Go logger
	// 	log.Fatal(err)
	// }

	// if err := e.StartTLS(":8443", "server.crt", "server.key"); err != http.ErrServerClosed {
	// 	log.Fatal(err)
	// }
}