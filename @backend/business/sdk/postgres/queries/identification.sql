-- name: InsertIdentification :exec
INSERT INTO identifications (
    id,
    first_name,
    last_name,
    middle_name,
    other_names,
    pin,
    identification_type,
    issued_country,
    place_of_birth,
    place_of_issue,
    date_of_birth,
    address_1,
    address_2,
    city,
    state_region,
    zip_code,
    created_at,
    updated_at
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,$17, $18)
RETURNING *;
