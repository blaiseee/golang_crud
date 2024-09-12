package handlers

import (
	"encoding/base64"
	"fmt"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"strings"

	"github.com/blaiseee/golang_crud/internal/models"
	"github.com/labstack/echo/v4"
)

func UploadImage(c echo.Context) error {
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
