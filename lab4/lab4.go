package main

import (
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Req struct {
	Message string `json:"message"`
}
type Word struct {
	Word   string `json:"word"`
	Length int    `json:"length"`
}

func main() {
	r := gin.Default()

	r.POST("/lab4", separateParagraph)
	r.GET("/gettime", getTimeEndpoint)
	r.POST("/postting", posttingEndpoint)

	r.Run(":8000")
}

func separateParagraph(c *gin.Context) {
	// request to struct
	var r Req
	c.ShouldBindJSON(&r)
	wordRes := []Word{}

	words := strings.Split(r.Message, " ")
	for _, w := range words {
		wordRes = append(wordRes, Word{w, len(w)})
	}

	// sort by length, word
	// Ref: https://golang.org/src/sort/example_test.go
	sort.SliceStable(wordRes, func(i, j int) bool {
		return wordRes[i].Word < wordRes[j].Word
	})
	sort.SliceStable(wordRes, func(i, j int) bool {
		return wordRes[i].Length < wordRes[j].Length
	})

	// return
	c.JSON(200, wordRes)
}

type Request struct {
	Time string
}

func getServerTime() time.Time {
	return time.Now()
}

func getTimeEndpoint(c *gin.Context) {
	t := getServerTime()
	trq := Request{}
	trq.Time = t.Format("02/01/2006 15:04:05")

	c.JSON(200, trq)
	// c.JSON(200, t.Format("Mon Jan 02-01-2006 15:04:05"))
	// c.JSON(200, t.Format(time.StampMicro))
}

func posttingEndpoint(c *gin.Context) {
	type Input struct {
		A int `json:"a"`
		B int `json:"b"`
	}
	type Response struct {
		A       int
		B       int
		Sum     int
		HourSum int
		Time    string
	}

	t := getServerTime()

	// bind JSON
	input := Input{}
	c.BindJSON(&input)

	// return JSON
	res := Response{}
	res.A = input.A
	res.B = input.B
	res.Sum = input.A + input.B
	res.HourSum = t.Day() + res.Sum
	res.Time = t.Format(time.ANSIC)
	c.JSON(200, res)
}
