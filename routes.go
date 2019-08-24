package main

import (
	"github.com/gin-gonic/gin"
)

// defineRoutes sets application REST endpoints
func (app *App) defineRoutes() {
	r := app.router

	// /ping app (http server is running)
	r.PUT("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.PUT("/hello/:username", putHelloUsername(app))
	r.GET("/hello/:username", getHelloUsername(app))
}
