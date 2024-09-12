package main

import (
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/blaiseee/golang_crud/internal/models"
	"github.com/blaiseee/golang_crud/internal/utils"
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

	e.POST("/upload-image", uploadImage)

	e.Logger.Fatal(e.Start(":8080"))
}

func getPosts(c echo.Context) error {
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

func getPost(c echo.Context) error {
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

	var output map[string]interface{}
	if err := utils.ParseJSONResponse(resp.Body(), &output); err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to parse JSON response")
	}

	return c.JSON(http.StatusCreated, output)
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

	var output map[string]interface{}
	if err := utils.ParseJSONResponse(resp.Body(), &output); err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to parse JSON response")
	}

	return c.JSON(http.StatusOK, output)
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

func uploadImage(c echo.Context) error {
	req := new(models.UploadImageRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid file format")
	}

	dataIndex := strings.Index(req.FileData, ",")
	if dataIndex != -1 {
		req.FileData = req.FileData[dataIndex+1:]
	}

	fileData, err := base64.StdEncoding.DecodeString(req.FileData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Failed to decode base64 string")
	}

	fileExtension := strings.ToLower(strings.Split(req.FileName, ".")[1])

	file, err := os.Create(req.FileName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to create image file")
	}

	defer file.Close()

	switch fileExtension {
	case "jpg", "jpeg":
		img, err := jpeg.Decode(strings.NewReader(string(fileData)))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Failed to decode JPEG image")
		}
		if err = jpeg.Encode(file, img, nil); err != nil {
			return c.JSON(http.StatusInternalServerError, "Failed to encode JPEG image")
		}
	case "png":
		img, err := png.Decode(strings.NewReader(string(fileData)))
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "Failed to decode PNG image")
		}
		if err = png.Encode(file, img); err != nil {
			return c.JSON(http.StatusInternalServerError, "Failed to encode PNG image")
		}
	default:
		return c.JSON(http.StatusInternalServerError, "Unsupported image file format")
	}
	return c.JSON(http.StatusOK, fmt.Sprintf("Image successfully saved as: %s", req.FileName))
}
