-- name: InsertOTP :exec
INSERT INTO otps(id, entity_id, hashed_code, code, expires_at, channel, destination)
VALUES($1, $2, $3, $4, $5, $6, $7);

-- name: GetOTP :one
SELECT id, entity_id, destination, channel, code, verified_at, expires_at, created_at, updated_at
FROM otps
WHERE id = $1;

-- name: GetOTPByCustomerIDAndCode :one
SELECT id, entity_id, hashed_code, code, verified_at, expires_at, created_at, updated_at
FROM otps
WHERE entity_id = $1
AND hashed_code = $2;
