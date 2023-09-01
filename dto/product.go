package dto

type ProductRequest struct {
	SKU        string  `json:"sku" binding:"required"`
	Name       string  `json:"name" binding:"required"`
	Desc       string  `json:"desc"`
	Price      float64 `json:"price" binding:"required"`
	Status     uint    `json:"status" binding:"required"`
	CategoryID uint    `json:"categoryId" binding:"required"`
	Image      []byte  `json:"image"`
}

type UpdateProductRequest struct {
	SKU        string  `json:"sku"`
	Name       string  `json:"name"`
	Desc       string  `json:"desc" `
	Price      float64 `json:"price" `
	Status     uint    `json:"status" `
	CategoryID uint    `json:"categoryId"`
}
type CreateOrUpdateProductResponse struct {
	ID         uint    `json:"id"`
	SKU        string  `json:"sku"`
	Name       string  `json:"name"`
	Desc       string  `json:"desc"`
	Price      float64 `json:"price"`
	Status     uint    `json:"status"`
	CategoryID uint    `json:"categoryId"`
	Image      string  `json:"image"`
}
type ReadProductResponse struct {
	ID       uint             `json:"id"`
	SKU      string           `json:"sku"`
	Name     string           `json:"name"`
	Desc     string           `json:"desc"`
	Price    float64          `json:"price"`
	Status   uint             `json:"status"`
	Category CategoryResponse `json:"category"`
	Image    string           `json:"image"`
}
