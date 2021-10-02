package service

import (
	"fmt"
	"math"

	"gorm.io/gorm"
)

const (
	defaultPageSize = 10
)

// PaginationOptions contains the query params associated with pagination
type PaginationOptions struct {
	Page     *int    `schema:"page" validate:"omitempty,min=1"`
	PageSize *int    `schema:"page_size" validate:"omitempty,min=1"`
	OrderBy  *string `schema:"order_by" validate:"omitempty,oneof=date_purchased"`
	Order    *string `schema:"order" validate:"omitempty,oneof=asc desc"`
}

// Paging contains the response on which page it is.
type Paging struct {
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Pages int   `json:"pages"`
}

// PaginatedResults is a generic struct that contains the response.
// Every list endpoint will contain a PaginatedResults struct.
type PaginatedResults struct {
	Results interface{} `json:"results"`
	Paging  Paging      `json:"paging"`
}

// PaginationScope contains the function to paginate the query.
func PaginationScope(options PaginationOptions) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if options.PageSize != nil {
			db = db.Limit(*options.PageSize)
			if options.Page != nil {
				db = db.Offset((*options.Page - 1) * *options.PageSize)
			}
		} else {
			db = db.Limit(defaultPageSize)
		}

		if options.OrderBy != nil {
			order := "asc"
			if options.Order != nil {
				order = *options.Order
			}
			db.Order(fmt.Sprintf("%s %s", *options.OrderBy, order))
		}

		return db
	}
}

// createPaginatedResults is a helper function that helps to create paginated results
func createPaginatedResults(options PaginationOptions, count int64, results interface{}) *PaginatedResults {
	page := 1
	totalPages := 1
	if options.Page != nil && options.PageSize != nil {
		page = int(math.Max(1, float64(*options.Page)))
		totalPages = int(math.Ceil(float64(count) / float64(*options.PageSize)))
	}
	return &PaginatedResults{
		Results: results,
		Paging: Paging{
			Total: count,
			Page:  page,
			Pages: totalPages,
		},
	}
}
