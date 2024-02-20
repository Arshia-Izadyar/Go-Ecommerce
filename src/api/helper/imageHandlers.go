package helper

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Arshia-Izadyar/Go-Ecommerce/src/api/dto"
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/pkg/service_errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SaveImages(ctx *gin.Context) (*dto.CreateCategoryDTO, error) {

	form, err := ctx.MultipartForm()
	if err != nil {
		return nil, &service_errors.ServiceError{EndUserMessage: "please send form data"}
	}
	var images []string

	files := form.File["images"]
	for i := 0; i < len(files); i++ {
		file := files[i]
		if strings.HasPrefix(file.Header["Content-Type"][0], "image/*") {
			return nil, &service_errors.ServiceError{EndUserMessage: "error please upload image/*"}

		}
		oldFileName := strings.Split(file.Filename, ".")
		file.Filename = fmt.Sprintf("%s.%s", uuid.New(), oldFileName[len(oldFileName)-1])
		dest := "../static/uploads/" + file.Filename
		err = ctx.SaveUploadedFile(file, dest)
		if err != nil {
			ctx.JSON(200, gin.H{"ok": err.Error()})
		}
		images = append(images, "/static/uploads/"+file.Filename)
	}

	req := &dto.CreateCategoryDTO{
		Name:   form.Value["name"][0],
		Slug:   form.Value["slug"][0],
		Images: images,
	}
	return req, nil
}

func RemoveImages(images []string) error {
	for i := 0; i < len(images); i++ {
		image := strings.Split(images[i], "/")

		imageName := image[len(image)-1]
		path := filepath.Join("../static/uploads/", imageName)
		err := os.Remove(path)
		if err != nil {
			return err
		}
	}
	return nil
}
