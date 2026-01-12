-- name: InsertBusiness :exec
INSERT INTO businesses (
    id,
    legal_name,
    registration_number,
    org_id,
    tax_id,
    entity,
    jurisdiction,
    dba,
    owners,
    address,
    website,
    phone_numbers,
    email_addresses,
    created_at,
    updated_at
) VALUES (@id,
@legal_name,
@registration_number,
@org_id,
@tax_id,
@entity,
@jurisdiction,
@dba,
@owners,
@address,
@website,
@phone_numbers,
@email_addresses,
@created_at,
@updated_at);


-- name: GetBusinessByID :one
SELECT * FROM businesses WHERE businesses.id = @id AND businesses.org_id = @org_id;

-- name: UpdateBusinessByID :exec
UPDATE businesses
SET
    legal_name      = @legal_name,
    tax_id          = @tax_id,
    entity          = @entity,
    jurisdiction    = @jurisdiction,
    dba             = @dba,
    address         = @address,   -- JSONB
    website         = @website,
    phone_numbers   = @phone_numbers,  -- ARRAY or JSONB
    email_addresses = @email_addresses,  -- ARRAY or JSONB
    updated_at = @updated_at,
    registration_number = @registration_number
WHERE id = @id AND org_id = @org_id;


-- name: DeleteByID :exec
DELETE FROM businesses WHERE businesses.id = $1 AND businesses.org_id = $2;
