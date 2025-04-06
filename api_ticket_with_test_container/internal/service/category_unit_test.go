package service

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github/moura95/ticket-api/config"
	"github/moura95/ticket-api/internal/repository"
)

func TestCategoryService_GetAll(t *testing.T) {

	tests := []struct {
		name          string
		parentId      string
		setupMocked   func(repo *repository.QuerierMocked)
		expected      []repository.Category
		expectedError error
	}{
		{
			name:     "Get sub categories",
			parentId: "1",
			setupMocked: func(repo *repository.QuerierMocked) {
				categories := []repository.Category{
					{
						ID:   3,
						Name: "Mouse",
						ParentID: sql.NullInt32{
							Int32: 1,
							Valid: true,
						},
					},
					{
						ID:   4,
						Name: "Keyboard",
						ParentID: sql.NullInt32{
							Int32: 1,
							Valid: true,
						},
					},
					{
						ID:   5,
						Name: "HeadPhone",
						ParentID: sql.NullInt32{
							Int32: 1,
							Valid: true,
						},
					},
				}
				repo.EXPECT().ListSubCategories(mock.Anything, mock.Anything).Return(categories, nil)

			},
			expected: []repository.Category{
				{
					ID:   3,
					Name: "Mouse",
					ParentID: sql.NullInt32{
						Int32: 1,
						Valid: true,
					},
				},
				{
					ID:   4,
					Name: "Keyboard",
					ParentID: sql.NullInt32{
						Int32: 1,
						Valid: true,
					},
				},
				{
					ID:   5,
					Name: "HeadPhone",
					ParentID: sql.NullInt32{
						Int32: 1,
						Valid: true,
					},
				},
			},
			expectedError: nil,
		},
		{
			name:     "Get all categories",
			parentId: "",
			setupMocked: func(repo *repository.QuerierMocked) {
				categories := []repository.Category{
					{
						ID:   1,
						Name: "Hardware",
					},
					{
						ID:   2,
						Name: "Software",
					},
				}
				repo.EXPECT().ListCategories(mock.Anything).Return(categories, nil)

			},
			expected: []repository.Category{
				{
					ID:   1,
					Name: "Hardware",
				},
				{
					ID:   2,
					Name: "Software",
				},
			},
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMocked := repository.NewQuerierMocked(t)
			tt.setupMocked(repoMocked)
			service := NewCategoryService(repoMocked, config.Config{}, nil)
			categories, err := service.GetAll(context.Background(), tt.parentId)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expected, categories)
		})
	}
}

func TestCategoryService_GetById(t *testing.T) {

	tests := []struct {
		name          string
		id            int32
		setupMocked   func(repo *repository.QuerierMocked)
		expected      *repository.Category
		expectedError error
	}{
		{
			name: "Get By Id Software Category",
			id:   1,
			setupMocked: func(repo *repository.QuerierMocked) {
				cate := repository.Category{
					ID:   1,
					Name: "Hardware",
				}
				repo.EXPECT().GetCategoryById(mock.Anything, mock.Anything).Return(cate, nil)

			},
			expected: &repository.Category{
				ID:   1,
				Name: "Hardware",
			},
		},
		{
			name: "Get By ID SubCategory",
			id:   4,
			setupMocked: func(repo *repository.QuerierMocked) {
				cate := repository.Category{
					ID:   4,
					Name: "Mouse",
				}
				repo.EXPECT().GetCategoryById(mock.Anything, mock.Anything).Return(cate, nil)

			},
			expected: &repository.Category{
				ID:   4,
				Name: "Mouse",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMocked := repository.NewQuerierMocked(t)
			tt.setupMocked(repoMocked)
			service := NewCategoryService(repoMocked, config.Config{}, nil)
			categories, err := service.GetByID(context.Background(), tt.id)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expected, categories)
		})
	}
}

func TestCategoryService_Create(t *testing.T) {

	tests := []struct {
		name          string
		categoryName  string
		parentId      int32
		setupMocked   func(repo *repository.QuerierMocked)
		expected      *repository.Category
		expectedError error
	}{
		{
			name: "Create Category Software",
			setupMocked: func(repo *repository.QuerierMocked) {
				cate := repository.Category{
					ID:   1,
					Name: "Software",
				}
				repo.EXPECT().CreateCategory(mock.Anything, mock.Anything).Return(cate, nil)

			},
			expected: &repository.Category{
				ID:   1,
				Name: "Software",
			},
		},
		{
			name: "Create Category Hardware",
			setupMocked: func(repo *repository.QuerierMocked) {
				cate := repository.Category{
					ID:   4,
					Name: "Hardware",
				}
				repo.EXPECT().CreateCategory(mock.Anything, mock.Anything).Return(cate, nil)

			},
			expected: &repository.Category{
				ID:   4,
				Name: "Hardware",
			},
		},
		{
			name:     "Create SubCategory Mouse in Hardware",
			parentId: 1,
			setupMocked: func(repo *repository.QuerierMocked) {
				cate := repository.Category{
					ID:   5,
					Name: "Mouse",
					ParentID: sql.NullInt32{
						Int32: 1,
						Valid: true,
					},
				}
				repo.EXPECT().CreateCategory(mock.Anything, mock.Anything).Return(cate, nil)

			},
			expected: &repository.Category{
				ID:   5,
				Name: "Mouse",
				ParentID: sql.NullInt32{
					Int32: 1,
					Valid: true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMocked := repository.NewQuerierMocked(t)
			tt.setupMocked(repoMocked)
			service := NewCategoryService(repoMocked, config.Config{}, nil)
			cate, err := service.Create(context.Background(), tt.categoryName, tt.parentId)
			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expected, cate)
		})
	}
}

func TestCategoryService_Delete(t *testing.T) {

	tests := []struct {
		name          string
		id            int32
		setupMocked   func(repo *repository.QuerierMocked)
		expectedError error
	}{
		{
			name: "Delete Category By Id 1",
			id:   1,
			setupMocked: func(repo *repository.QuerierMocked) {
				repo.EXPECT().DeleteCategory(mock.Anything, mock.Anything).Return(nil)

			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMocked := repository.NewQuerierMocked(t)
			tt.setupMocked(repoMocked)
			service := NewCategoryService(repoMocked, config.Config{}, nil)
			err := service.Delete(context.Background(), tt.id)
			assert.Equal(t, tt.expectedError, err)

		})
	}
}
