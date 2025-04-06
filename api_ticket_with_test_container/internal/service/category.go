package service

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github/moura95/ticket-api/config"
	"github/moura95/ticket-api/internal/repository"
	"go.uber.org/zap"
)

type CategoryService struct {
	repository repository.Querier
	config     config.Config
	logger     *zap.SugaredLogger
}

func NewCategoryService(repo repository.Querier, cfg config.Config, log *zap.SugaredLogger) *CategoryService {
	return &CategoryService{
		repository: repo,
		config:     cfg,
		logger:     log,
	}
}

func (s *CategoryService) GetAll(ctx context.Context, parentId string) (categories []repository.Category, error error) {
	if parentId != "" {
		// Conver string to int32
		id, err := strconv.Atoi(parentId)
		if err != nil {
			return nil, fmt.Errorf("failed to convert id %s", err.Error())
		}
		categories, err = s.repository.ListSubCategories(ctx, sql.NullInt32{
			Int32: int32(id),
			Valid: id > 0,
		})
		return categories, nil

	}
	categories, err := s.repository.ListCategories(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get category %s", err.Error())
	}

	return categories, nil
}

func (s *CategoryService) GetByID(ctx context.Context, id int32) (*repository.Category, error) {
	category, err := s.repository.GetCategoryById(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("failed to get category %s", err.Error())
	}

	return &category, nil
}

func (s *CategoryService) Create(ctx context.Context, name string, parentId int32) (*repository.Category, error) {
	arg := repository.CreateCategoryParams{
		Name: name,
		ParentID: sql.NullInt32{
			Int32: parentId,
			Valid: parentId > 0,
		},
	}
	category, err := s.repository.CreateCategory(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("failed to create %s", err.Error())
	}
	return &category, nil
}

func (s *CategoryService) Delete(ctx context.Context, id int32) error {
	err := s.repository.DeleteCategory(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get category %s", err.Error())
	}

	return nil
}

func (s *CategoryService) Update(ctx context.Context, id, parentID int32, name string) error {
	arg := repository.UpdateCategoryParams{
		ID: id,
		Name: sql.NullString{
			String: name,
			Valid:  name != "",
		},
		ParentID: sql.NullInt32{
			Int32: parentID,
			Valid: parentID > 0,
		},
	}
	err := s.repository.UpdateCategory(ctx, arg)
	if err != nil {
		return fmt.Errorf("failed to update %s", err.Error())
	}
	return nil
}
