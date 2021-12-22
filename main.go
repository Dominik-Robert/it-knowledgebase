package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// test
func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/**/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	router.GET("/detail/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "detail.html", gin.H{})
	})

	router.StaticFS("/assets", http.Dir("assets/"))

	router.Run(":8080")
}
