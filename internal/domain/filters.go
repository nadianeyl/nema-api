package domain

import (
	"math"
	"strings"
)

type Filters struct {
	Limit        int
	Page         int
	Sort         string
	SortSafelist []string
}

type Metadata struct {
	Limit        int `json:"limit,omitempty"`
	CurrentPage  int `json:"current_page,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

func (f Filters) SortColumn() string {
	for _, safeValue := range f.SortSafelist {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}

	panic("unsafe sort parameter: " + f.Sort)
}

func (f Filters) SortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}

	return "ASC"
}

func (f Filters) GetLimit() int {
	return f.Limit
}

func (f Filters) GetPage() int {
	return f.Page
}

func (f Filters) GetOffset() int {
	return (f.Page - 1) * f.Limit
}

func GenerateMetadata(limit, page, totalRecords int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		Limit:        limit,
		CurrentPage:  page,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(limit))),
		TotalRecords: totalRecords,
	}
}
