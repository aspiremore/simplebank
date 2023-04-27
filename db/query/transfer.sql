-- name: CreateTransfer :one
INSERT INTO transfers (
    from_account_id,
    to_account_id,
    amount
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: UpdateTransfer :one
Update transfers
SET amount = $3
WHERE from_account_id = $1 and to_account_id = $2
    RETURNING *;


-- name: GetTransfer :one
SELECT * FROM transfers
WHERE ID = $1 LIMIT 1;