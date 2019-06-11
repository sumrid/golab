package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	basicEcho()
	// basicGroup()
}

func basicEcho() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
	// 	if username == "joe" && password == "secret" {
	// 		return true, nil
	// 	}
	// 	return false, nil
	// }))

	e.GET("/", index).Name = "index"
	e.GET("/grade", gradeEndpoint)
	e.POST("/book", insertBook)
	e.GET("/book", getBook)
	e.POST("/", postEndpoint)

	e.Logger.Info(e.Routes())
	e.Logger.Fatal(e.Start(":80"))
}

func basicGroup() {
	e := echo.New()
	gr := e.Group("/v1")
	gr.GET("/", index)
	e.Logger.Fatal(e.Start(":80"))
}

// Req is Struct for request
type Req struct {
	Name string
}

// Res is stuct for response
type Res struct {
	Time     string
	Messsage string
	Name     string
}

func index(c echo.Context) error {
	name := c.QueryParam("name")
	return c.JSON(http.StatusOK, Res{time.Now().Format(time.ANSIC), "hello", name})
}

func postEndpoint(c echo.Context) error {
	// Struct for req
	type Req struct {
		Name string
	}

	// req := Req{}
	req := new(Req) // return pointer of struct
	err := c.Bind(req)
	if err != nil {
		return err
	}

	res := Res{}
	res.Time = time.Now().Format(time.ANSIC)
	res.Name = req.Name
	return c.JSON(http.StatusOK, res)
}

func gradeEndpoint(c echo.Context) error {
	p := c.QueryParam("point")
	pInt, err := strconv.Atoi(p)
	grade := ""

	if err != nil {
		grade = "Input is not number"
	} else if pInt <= 100 && pInt >= 0 {
		if pInt > 90 {
			grade = "A"
		} else if pInt > 80 {
			grade = "B"
		} else if pInt > 70 {
			grade = "C"
		} else if pInt > 60 {
			grade = "D"
		} else {
			grade = "F"
		}
	} else {
		grade = "Input invalid!! point 0-100"
	}

	type ResPoint struct {
		Point int
		Grade string
	}

	return c.JSON(http.StatusOK, ResPoint{pInt, grade})
}
