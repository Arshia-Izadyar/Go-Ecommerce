package dto

type CreateCategoryDTO struct {
	Name   string
	Slug   string
	Images []string
}

type CategoryResponse struct {
	Id     int      `json:"id"`
	Name   string   `json:"name"`
	Slug   string   `json:"slug"`
	Images []string `json:"images"`
}

type UpdateCategoryDTO struct {
	Name   string   `json :""`
	Images []string `json :""`
	Append bool     `json :""`
}
