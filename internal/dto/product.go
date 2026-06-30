package dto

type CreateProductRequest struct {
	Code        string `json:"code" binding:"required,max=64"`
	Name        string `json:"name" binding:"required,max=255"`
	BillingType string `json:"billing_type" binding:"required,oneof=monthly one_time"`
	Price       int64  `json:"price" binding:"required,gte=0"`
	IsActive    *bool  `json:"is_active"`
}

type UpdateProductRequest struct {
	Code        *string `json:"code" binding:"omitempty,max=64"`
	Name        *string `json:"name" binding:"omitempty,max=255"`
	BillingType *string `json:"billing_type" binding:"omitempty,oneof=monthly one_time"`
	Price       *int64  `json:"price" binding:"omitempty,gte=0"`
	IsActive    *bool   `json:"is_active"`
}

type ProductResponse struct {
	ID          uint   `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	BillingType string `json:"billing_type"`
	Price       int64  `json:"price"`
	IsActive    bool   `json:"is_active"`
}

type PaginatedProductsResponse struct {
	Data       []ProductResponse `json:"data"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalItems int64             `json:"total_items"`
	TotalPages int64             `json:"total_pages"`
}
