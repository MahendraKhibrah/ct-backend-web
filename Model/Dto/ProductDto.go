package Dto

type AddProductRequest struct {
	Name string `json:"name" binding:"required"`
}

type EditProductRequest struct {
	Name string `json:"name" binding:"required"`
	Id   int    `json:"id" binding:"required"`
}

type GetProductRequest struct {
	Search string `form:"search"`
}
