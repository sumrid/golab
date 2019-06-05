package main

import "github.com/gin-gonic/gin"
import "net/http"
import "strconv"

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello world")
	})

	r.GET("/lap3", func(c *gin.Context) {
		point := c.Query("point")
		pointInt, err := strconv.Atoi(point)
		grade := ""

		if err == nil && pointInt <= 100 && pointInt >= 0 {
			if pointInt > 90 {
				grade = "A"
			} else if pointInt > 80 {
				grade = "B"
			} else if pointInt > 70 {
				grade = "C"
			} else if pointInt > 60 {
				grade = "D"
			} else {
				grade = "F"
			}
		} else {
			grade = "Input invalid!!"
		}

		c.JSON(http.StatusOK, gin.H{
			"point": point,
			"grade": grade,
		})
	})

	r.Run(":8000")
}
