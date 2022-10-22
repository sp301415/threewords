-- name: CheckID :one
SELECT * FROM fileDB WHERE id = ?;

-- name: CreateEntry :exec
INSERT INTO fileDB (id, file_path, original_name, expire_date) VALUES (?, ?, ?, unixepoch('now', '+24 hours'));

-- name: ReadEntry :one
SELECT file_path, original_name FROM fileDB WHERE id = ?;

-- name: FindExpiredEntry :many
SELECT id, file_path FROM fileDB WHERE expire_date < unixepoch('now');

-- name: DeleteEntry :exec
DELETE FROM fileDB WHERE id = ?;
