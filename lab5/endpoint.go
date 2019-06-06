package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func pingEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func pingV2Endpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong v2"})
}

func helloEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "hello world"})
}
