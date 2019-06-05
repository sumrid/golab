package main

import (
	"sort"
	"strings"

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
