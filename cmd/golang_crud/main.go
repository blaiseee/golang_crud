package main

import (
	"net/http"
	"strconv"

	"github.com/blaiseee/golang_crud/internal/models"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
)

const baseURL = "https://jsonplaceholder.typicode.com"

func main() {
	e := echo.New()

	e.GET("/posts", getPosts)
	e.GET("/posts/:id", getPost)
	e.POST("/posts", createPost)
	e.PUT("/posts/:id", updatePost)
	e.DELETE("/posts/:id", deletePost)

	e.Logger.Fatal(e.Start(":8080"))
}

func getPosts(c echo.Context) error {
	client := resty.New()
	resp, err := client.R().Get(baseURL + "/posts")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, resp.String())
}

func getPost(c echo.Context) error {
	id := c.Param("id")
	client := resty.New()
	resp, err := client.R().Get(baseURL + "/posts/" + id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, resp.String())
}

func createPost(c echo.Context) error {
	post := new(models.Post)
	if err := c.Bind(post); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	client := resty.New()
	resp, err := client.R().
		SetBody(post).
		Post(baseURL + "/posts")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, resp.String())
}

func updatePost(c echo.Context) error {
	id := c.Param("id")
	post := new(models.Post)
	if err := c.Bind(post); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	client := resty.New()
	resp, err := client.R().
		SetBody(post).
		Put(baseURL + "/posts/" + id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, resp.String())
}

func deletePost(c echo.Context) error {
	id := c.Param("id")
	client := resty.New()
	resp, err := client.R().Delete(baseURL + "/posts/" + id)
	if err != nil {
		c.Logger().Errorf("Error deleting post: %v", err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "Post deleted successfully: "+strconv.Itoa(resp.StatusCode()))
}
