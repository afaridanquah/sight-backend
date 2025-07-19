-- name: CreateCustomer :exec
    INSERT INTO customers(id, first_name, last_name, middle_name, date_of_birth, birth_country, city_of_birth, email, phone_number, business_id, creator_id, identifications, addresses, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
    RETURNING *;

-- name: QueryCustomerByID :one
    SELECT * FROM customers
    WHERE customers.id = $1;

-- name: QueryCustomerByAndBusinessID :one
    SELECT * FROM customers
    WHERE customers.id = $1
    AND customers.business_id = $2;
