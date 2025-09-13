-- name: CreateUser :one
INSERT INTO users (email, password_hash)
VALUES ($1, $2)
RETURNING id, created_at;

-- name: GetUserByEmail :one
SELECT id, password_hash, email
FROM users
WHERE email = $1;

-- name: CreateService :one
INSERT INTO services (user_id, name, target, check_interval_seconds)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetServicesAndOwners :many
SELECT s.*, u.email as owner_email
FROM services s
JOIN users u ON s.user_id = u.id;

-- name: DeleteService :execrows
DELETE FROM services
WHERE id = $1 AND user_id = $2;

-- name: CreateStatusCheck :one
INSERT INTO status_checks (service_id, status, status_code, response_time_ms, error_message)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetStatusChecksForService :many
SELECT sc.*
FROM status_checks sc
JOIN services s ON sc.service_id = s.id
WHERE sc.service_id = $1 AND s.user_id = $2
ORDER BY sc.checked_at DESC
LIMIT 50;

-- name: GetLatestStatusCheckForService :one
SELECT status FROM status_checks
WHERE service_id = $1
ORDER BY checked_at DESC
LIMIT 1;
