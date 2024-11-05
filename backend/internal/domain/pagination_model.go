// internal/domain/pagination.go
package domain

type PaginationQuery struct {
    Page     int `json:"page"`
    PageSize int `json:"page_size"`
}

type PaginatedResponse struct {
    Data       interface{} `json:"data"`
    Total      int64       `json:"total"`
    Page       int         `json:"page"`
    PageSize   int         `json:"page_size"`
    TotalPages int         `json:"total_pages"`
}

