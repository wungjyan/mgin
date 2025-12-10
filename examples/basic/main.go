package main

import (
	"mgin"
	"net/http"
)

func main() {
	r := mgin.New()

	r.GET("/", func(c *mgin.Context) {
		c.JSON(http.StatusOK, mgin.H{"message": "Hello, World!"})
	})

	r.POST("/login", func(c *mgin.Context) {
		c.JSON(http.StatusOK, mgin.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	r.Run(":8080")
}
