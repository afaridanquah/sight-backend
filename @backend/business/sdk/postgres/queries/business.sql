-- name: InsertBusiness :exec
INSERT INTO businesses (
    id,
    legal_name,
    org_id,
    tax_id,
    entity,
    jurisdiction,
    dba,
    admin_id,
    owners,
    address,
    website,
    phone_numbers,
    email_addresses,
    created_at,
    updated_at
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15);


-- name: GetBusinessByID :one
SELECT * FROM businesses WHERE id = $1 AND org_id = $2;

-- name: UpdateBusinessByID :exec
UPDATE businesses
SET
    legal_name      = $2,
    tax_id          = $3,
    entity          = $4,
    jurisdiction    = $5,
    dba             = $6,
    admin_id        = $7,
    address         = $8,   -- JSONB
    website         = $9,
    phone_numbers   = $10,  -- ARRAY or JSONB
    email_addresses = $11,  -- ARRAY or JSONB
    updated_at = $12
WHERE id = $1 AND org_id = $13;
