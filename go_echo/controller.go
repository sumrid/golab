package main

import (
	"time"

	"github.com/labstack/echo"
	"github.com/sumrid/golab/go_echo/model"
	"gopkg.in/mgo.v2"
)

func insertBook(c echo.Context) error {
	// Bind request from json
	bookReq := model.Book{}
	if err := c.Bind(&bookReq); err != nil {
		return err
	}
	bookReq.PublishDate = time.Now()

	// Create connection
	ss, err := mgo.Dial("mongodb://127.0.0.1:27017")
	if err != nil {
		return err
	}
	defer ss.Close()
	// Set mode
	ss.SetMode(mgo.Monotonic, true)

	// Select DB and collection
	col := ss.DB("mydb").C("books")
	if err := col.Insert(bookReq); err != nil {
		return err
	}

	return c.JSON(200, bookReq)
}

func getBook(c echo.Context) error {
	ss, err := mgo.Dial("mongodb://127.0.0.1:27017")
	if err != nil {
		return err
	}
	defer ss.Close()
	ss.SetMode(mgo.Monotonic, true)

	// Get all books
	result := []model.Book{}
	col := ss.DB("mydb").C("books")
	col.Find(nil).All(&result)
	return c.JSON(200, result)
}
