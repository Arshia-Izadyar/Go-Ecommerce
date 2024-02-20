package handler

import (
	"fmt"

	"github.com/Arshia-Izadyar/Go-Ecommerce/src/api/helper"
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/config"
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/services"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service *services.CategoryService
}

func NewCategoryHandler(cfg *config.Config) *CategoryHandler {
	s := services.NewCategoryService(cfg)
	return &CategoryHandler{
		service: s,
	}
}

func (ch *CategoryHandler) CreateCategory(ctx *gin.Context) {
	req, err := helper.SaveImages(ctx)
	if err != nil {
		ctx.AbortWithStatusJSON(400, helper.GenerateResponseWithError(err, 400, false))
		return
	}
	res, err := ch.service.CreateCategory(req)
	if err != nil {
		err = helper.RemoveImages(req.Images)
		if err != nil {
			fmt.Println(err)
		}
		ctx.AbortWithStatusJSON(400, helper.GenerateResponseWithError(err, 400, false))
	}
	ctx.JSON(200, gin.H{"ok": res})
}
