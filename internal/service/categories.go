package service

import (
	"github.com/nadianeyl/nema-api/internal/domain"
	"github.com/nadianeyl/nema-api/internal/repository"
)

type CategoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return CategoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *CategoryService) List(req *ListCategoriesRequest) ([]*CategoryResponse, domain.Metadata, error) {
	categories, metadata, err := s.categoryRepo.GetAllForUser(req.UserID, req.TransactionType, req.Filters)
	if err != nil {
		return nil, domain.Metadata{}, err
	}

	res := make([]*CategoryResponse, 0)
	for _, category := range categories {
		res = append(res, &CategoryResponse{
			ID:              category.ID,
			UserID:          category.UserID,
			Name:            category.Name,
			TransactionType: category.TransactionType,
			CreatedAt:       category.CreatedAt,
			UpdatedAt:       category.UpdatedAt,
		})
	}

	return res, metadata, nil
}
