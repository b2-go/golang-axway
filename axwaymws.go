package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	)

func setupRouter() *gin.Engine {
	var m map[string]string
	m = make(map[string]string)
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/v1/keyvalue/:key", func(c *gin.Context) {
		key := c.Param("key")
		value, ok := m[key]
		if (ok) {
			c.String(http.StatusOK, value)
		} else {
			c.String(http.StatusNotFound, "not found")
		}
	})
	r.POST("/v1/keyvalue", func(c *gin.Context) {
		key := c.Query("key")
		value := c.Query("value")
		m[key] = value
		c.String(http.StatusOK, "ok: value set")
	})
	return r
}

func main() {
	r := setupRouter()
	r.Run(":8888") // listen and serve on 0.0.0.0:8888
}