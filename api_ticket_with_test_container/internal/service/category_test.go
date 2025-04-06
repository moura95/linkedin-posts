package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github/moura95/ticket-api/db"
	"github/moura95/ticket-api/internal/repository"
)

func NewCategoryServiceTest(repo repository.Querier) *CategoryService {
	return &CategoryService{
		repository: repo,
	}
}

type createCategoryRequest struct {
	Name     string `json:"name"`
	ParentId int32  `json:"parent_id"`
}

func TestGetCategory(t *testing.T) {
	ctx := context.Background()
	conn, err := db.ConnectPostgres(connStr)
	store := repository.New(conn.DB())

	service := NewCategoryServiceTest(store)

	category, err := service.GetByID(ctx, 1)

	assert.NoError(t, err)
	assert.Equal(t, "Hardware", category.Name)

}

func TestListCategories(t *testing.T) {
	ctx := context.Background()
	conn, err := db.ConnectPostgres(connStr)
	store := repository.New(conn.DB())

	service := NewCategoryServiceTest(store)

	categories, err := service.GetAll(ctx, "")

	assert.NoError(t, err)
	assert.Equal(t, len(categories), 12)
	assert.Equal(t, "Hardware", categories[0].Name)

}

func TestListSubCategories(t *testing.T) {
	ctx := context.Background()
	conn, err := db.ConnectPostgres(connStr)
	store := repository.New(conn.DB())

	service := NewCategoryServiceTest(store)

	categories, err := service.GetAll(ctx, "1")

	assert.NoError(t, err)
	assert.Equal(t, len(categories), 3)
	assert.Equal(t, "Laptops", categories[0].Name)

}
func TestCreateCategory(t *testing.T) {
	ctx := context.Background()
	conn, err := db.ConnectPostgres(connStr)
	store := repository.New(conn.DB())

	service := NewCategoryServiceTest(store)
	req := createCategoryRequest{
		Name: "Hardware",
	}

	ca, err := service.Create(ctx, req.Name, 0)
	assert.NoError(t, err)

	category, _ := service.GetByID(ctx, ca.ID)
	assert.Equal(t, "Hardware", category.Name)
	assert.Equal(t, false, category.ParentID.Valid)

	assert.NoError(t, err)

}

func TestCreateSubCategory(t *testing.T) {
	ctx := context.Background()
	conn, err := db.ConnectPostgres(connStr)
	store := repository.New(conn.DB())

	service := NewCategoryServiceTest(store)
	req := createCategoryRequest{
		Name:     "Computer",
		ParentId: 1,
	}

	ca, err := service.Create(ctx, req.Name, req.ParentId)
	assert.NoError(t, err)

	category, _ := service.GetByID(ctx, ca.ID)
	assert.Equal(t, "Computer", category.Name)
	assert.Equal(t, int32(1), category.ParentID.Int32)

	assert.NoError(t, err)

}

func TestDeleteCategory(t *testing.T) {
	ctx := context.Background()
	conn, err := db.ConnectPostgres(connStr)
	store := repository.New(conn.DB())

	service := NewCategoryServiceTest(store)

	// delete item
	err = service.Delete(ctx, 1)
	assert.NoError(t, err)

	// confirm list
	newCategoriesList, err := service.GetAll(ctx, "")
	assert.NoError(t, err)
	assert.Equal(t, 9, len(newCategoriesList))

}
