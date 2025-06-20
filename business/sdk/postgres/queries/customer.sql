-- name: CreateCustomer :one
    INSERT INTO customers(id, first_name, last_name, middle_name, email, country, business_id, user_id, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    RETURNING *;

-- name: FindCustomer :one
    SELECT * FROM customers
    INNER JOIN customer_identifications ON (customers.id = customer_identifications.customer_id)
    WHERE customers.id = $1;

-- name: CreateCustomerIdentifications :copyfrom
    INSERT INTO customer_identifications(
        id,
        issued_country,
        customer_id,
        identification_type,
        pin
    ) VALUES($1, $2, $3, $4, $5);

-- name: BulkInsertCustomerAddress :copyfrom
    INSERT INTO customer_addresses(
        id,
        address_1,
        address_2,
        city,
        state,
        zip,
        country
    ) VALUES($1, $2, $3, $4, $5, $6, $7);
