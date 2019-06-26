package main

import (
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
