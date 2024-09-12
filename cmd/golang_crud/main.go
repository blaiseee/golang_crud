package main

import (
	"github.com/blaiseee/golang_crud/handlers"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.GET("/posts", handlers.GetPosts)
	e.GET("/posts/:id", handlers.GetPost)
	e.POST("/posts", handlers.CreatePost)
	e.PUT("/posts/:id", handlers.UpdatePost)
	e.DELETE("/posts/:id", handlers.DeletePost)

	e.POST("/upload-image", handlers.UploadImage)

	e.Logger.Fatal(e.Start(":8080"))
}
