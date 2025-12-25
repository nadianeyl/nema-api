package domain

import "math"

type Filters struct {
	Limit int
	Page  int
}

type Metadata struct {
	Limit        int `json:"limit,omitempty"`
	CurrentPage  int `json:"current_page,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
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
