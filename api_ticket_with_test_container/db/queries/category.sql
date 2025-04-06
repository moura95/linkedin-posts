-- name: CreateCategory :one
INSERT INTO categories (name,parent_id)
VALUES ($1, $2)
RETURNING *;


-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = $1;

-- name: GetCategoryById :one
SELECT *
FROM categories
WHERE categories.id = $1;

-- name: ListCategories :many
SELECT id, name, parent_id
FROM categories;

-- name: ListSubCategories :many
SELECT id, name, parent_id
FROM categories
WHERE  parent_id = $1;



-- name: UpdateCategory :exec
UPDATE categories
SET name = COALESCE(sqlc.narg('name'), name),
    parent_id  = COALESCE(sqlc.narg('parent_id'), parent_id)
WHERE id = $1;