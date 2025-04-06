package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github/moura95/ticket-api/db"
)

func TestCategoryRepository_Create(t *testing.T) {
	ctx := context.Background()
	conn, err := db.ConnectPostgres(connStr)
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	store := New(conn.DB())

	arg := CreateCategoryParams{
		Name: "Category",
		ParentID: sql.NullInt32{
			Int32: 0,
			Valid: false,
		},
	}

	category, err := store.CreateCategory(ctx, arg)
	assert.NoError(t, err)
	assert.Equal(t, "Category", category.Name)

}

func TestSubCategoryRepository_Create(t *testing.T) {
	ctx := context.Background()
	conn, err := db.ConnectPostgres(connStr)
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	store := New(conn.DB())

	arg := CreateCategoryParams{
		Name: "Computer",
		ParentID: sql.NullInt32{
			Int32: 1,
			Valid: true,
		},
	}

	category, err := store.CreateCategory(ctx, arg)
	assert.NoError(t, err)
	assert.Equal(t, "Computer", category.Name)
	assert.Equal(t, int32(1), category.ParentID.Int32)

}

func TestCategoryRepository_Create_Failed(t *testing.T) {
	ctx := context.Background()
	conn, err := db.ConnectPostgres(connStr)
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	store := New(conn.DB())

	arg := CreateCategoryParams{
		Name: "Category Fail",
		ParentID: sql.NullInt32{
			Int32: 999999,
			Valid: true,
		},
	}

	_, err = store.CreateCategory(ctx, arg)
	assert.Equal(t, "pq: insert or update on table \"categories\" violates foreign key constraint \"categories_parent_id_fkey\"", err.Error())

}

func TestCategoryRepository_GetById(t *testing.T) {
	ctx := context.Background()
	conn, err := db.ConnectPostgres(connStr)
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	store := New(conn.DB())
	id := 1
	category, err := store.GetCategoryById(ctx, int32(id))
	assert.NoError(t, err)

	assert.Equal(t, "Hardware", category.Name)

}

func TestCategoryRepository_GetAll(t *testing.T) {
	ctx := context.Background()
	conn, err := db.ConnectPostgres(connStr)
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	store := New(conn.DB())
	categories, err := store.ListCategories(ctx)
	assert.NoError(t, err)

	// category
	assert.Equal(t, "Hardware", categories[0].Name)

	// category 2
	assert.Equal(t, "Software", categories[1].Name)

}

func TestCategoryRepository_UpdateName(t *testing.T) {
	ctx := context.Background()
	conn, err := db.ConnectPostgres(connStr)
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	store := New(conn.DB())

	arg := UpdateCategoryParams{
		ID: 1,
		Name: sql.NullString{
			String: "Update",
			Valid:  true,
		},
		ParentID: sql.NullInt32{},
	}

	err = store.UpdateCategory(ctx, arg)
	assert.NoError(t, err)

	// Get Category for confirm update

	cate, err := store.GetCategoryById(ctx, arg.ID)
	assert.Equal(t, "Update", cate.Name)

}

func TestCategoryRepository_UpdateNameAndParent(t *testing.T) {
	ctx := context.Background()
	conn, err := db.ConnectPostgres(connStr)
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	store := New(conn.DB())

	arg := UpdateCategoryParams{
		ID: 1,
		Name: sql.NullString{
			String: "Update New",
			Valid:  true,
		},
		ParentID: sql.NullInt32{
			Int32: 1,
			Valid: true,
		},
	}

	err = store.UpdateCategory(ctx, arg)
	assert.NoError(t, err)

	// Get Category for confirm update

	cate, err := store.GetCategoryById(ctx, arg.ID)
	assert.Equal(t, "Update New", cate.Name)
	assert.Equal(t, int32(1), cate.ParentID.Int32)

}

func TestCategoryRepository_Delete(t *testing.T) {
	ctx := context.Background()
	conn, err := db.ConnectPostgres(connStr)
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	store := New(conn.DB())
	id := 1
	err = store.DeleteCategory(ctx, int32(id))
	assert.NoError(t, err)

}
