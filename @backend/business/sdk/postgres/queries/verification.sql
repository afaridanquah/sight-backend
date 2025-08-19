-- name: InsertVerification :exec
INSERT INTO verifications (id, verification_type, customer_id, customer, business_id, creator_id, outcome, aml_insight, phone_insight, decision, created_at, updated_at)
VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);
