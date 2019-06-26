package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/sumrid/golab/go_echo/model"
	"gopkg.in/mgo.v2"
)

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

func insertBook(c echo.Context) error {
	// นำเอา request ไปใส่ใน struct
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
