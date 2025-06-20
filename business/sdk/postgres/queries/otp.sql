-- name: InsertOTP :one
INSERT INTO otps(id, customer_id, hashed_code, expires_at, channel, destination)
VALUES($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetOTP :one
SELECT id, customer_id, destination, channel, hashed_code, verified_at, expires_at, created_at, updated_at
FROM otps
WHERE id = $1;

-- name: GetOTPByCustomerIDAndHash :one
SELECT id, customer_id, destination, channel, hashed_code, verified_at, expires_at, created_at, updated_at
FROM otps
WHERE customer_id = $1
AND hashed_code = $2;


-- name: UpdateOTP :exec
UPDATE otps 
SET verified_at = NOW() 
WHERE customer_id = $1 AND hashed_code = $2;