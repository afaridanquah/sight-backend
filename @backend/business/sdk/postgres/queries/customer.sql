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

-- name: UpdateCustomer :exec
    UPDATE customers
    SET first_name = $2,
        last_name = $3,
        middle_name = $4,
        birth_country = $5,
        city_of_birth = $6,
        email = $7,
        phone_number = $8,
        updated_at = $9,
        date_of_birth = $10
    WHERE id = $1;
