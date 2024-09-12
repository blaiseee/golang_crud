package handlers

import (
	"net/http"
	"strconv"

	"github.com/blaiseee/golang_crud/internal/models"
	"github.com/blaiseee/golang_crud/utils"
	"github.com/go-resty/resty/v2"
	"github.com/labstack/echo/v4"
)

const baseURL = "https://jsonplaceholder.typicode.com"

func GetPosts(c echo.Context) error {
	client := resty.New()
	resp, err := client.R().Get(baseURL + "/posts")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	var posts []map[string]interface{}
	if err := utils.ParseJSONResponse(resp.Body(), &posts); err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to parse JSON response")
	}

	return c.JSON(http.StatusOK, posts)
}

func GetPost(c echo.Context) error {
	id := c.Param("id")
	client := resty.New()
	resp, err := client.R().Get(baseURL + "/posts/" + id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	var post map[string]interface{}
	if err := utils.ParseJSONResponse(resp.Body(), &post); err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to parse JSON response")
	}

	return c.JSON(http.StatusOK, post)
}

func CreatePost(c echo.Context) error {
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

	var output map[string]interface{}
	if err := utils.ParseJSONResponse(resp.Body(), &output); err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to parse JSON response")
	}

	return c.JSON(http.StatusCreated, output)
}

func UpdatePost(c echo.Context) error {
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

	var output map[string]interface{}
	if err := utils.ParseJSONResponse(resp.Body(), &output); err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to parse JSON response")
	}

	return c.JSON(http.StatusOK, output)
}

func DeletePost(c echo.Context) error {
	id := c.Param("id")
	client := resty.New()
	resp, err := client.R().Delete(baseURL + "/posts/" + id)
	if err != nil {
		c.Logger().Errorf("Error deleting post: %v", err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "Post deleted successfully: "+strconv.Itoa(resp.StatusCode()))
}
