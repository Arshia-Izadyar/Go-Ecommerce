package services

import (
	"encoding/json"
	"fmt"

	"github.com/Arshia-Izadyar/Go-Ecommerce/src/api/dto"
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/config"
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/data/database"
	"github.com/Arshia-Izadyar/Go-Ecommerce/src/data/models"
	"gorm.io/gorm"
)

type CategoryService struct {
	Db  *gorm.DB
	Cfg *config.Config
}

func NewCategoryService(cfg *config.Config) *CategoryService {
	db := database.GetDB()
	return &CategoryService{
		Db:  db,
		Cfg: cfg,
	}
}

func (cs *CategoryService) CreateCategory(req *dto.CreateCategoryDTO) (*dto.CategoryResponse, error) {
	tx := cs.Db.Begin()

	res, err := json.Marshal(req.Images)
	if err != nil {
		return nil, err
	}
	// req.Images = string(req)
	category := &models.Category{
		Name:   req.Name,
		Slug:   req.Slug,
		Images: []string{string(res)},
	}

	err = tx.Model(&models.Category{}).Create(category).Error
	if err != nil {
		fmt.Println(err)
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return nil, nil
}
