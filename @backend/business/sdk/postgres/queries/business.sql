-- name: InsertBusiness :exec
INSERT INTO businesses (
    id,
    legal_name,
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
) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14);


-- name: GetBusinessByID :one
SELECT * FROM businesses WHERE id = $1;

-- name: GetBusinessByTaxID :one
SELECT * FROM businesses WHERE tax_id = $1;
