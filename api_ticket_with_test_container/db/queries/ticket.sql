-- name: CreateTicket :one
INSERT INTO tickets (title,description,severity_id,category_id,subcategory_id,status)
VALUES ($1, $2, $3, $4, $5, 'OPEN')
RETURNING *;


-- name: DeleteTicket :exec
DELETE FROM tickets
WHERE id = $1;

-- name: GetTicketById :one
SELECT tickets.*, category.name as category_name, subcategory.name as subcategory_name
FROM tickets
JOIN categories category on tickets.category_id = category.id
LEFT JOIN categories subcategory on tickets.subcategory_id = subcategory.id
WHERE tickets.id = $1;

-- name: ListTickets :many
SELECT *
FROM tickets
order by id;

-- name: UpdateTicket :exec
UPDATE tickets
SET title = COALESCE(sqlc.narg('title'), title),
    description  = COALESCE(sqlc.narg('description'), description),
    user_id      = COALESCE(sqlc.narg('user_id'), user_id),
    severity_id      = COALESCE(sqlc.narg('severity_id'), severity_id),
    category_id      = COALESCE(sqlc.narg('category_id'), category_id),
    subcategory_id      = COALESCE(sqlc.narg('subcategory_id'), subcategory_id),
    status        = COALESCE(sqlc.narg('status'), status),
    updated_at = NOW()
WHERE id = $1;