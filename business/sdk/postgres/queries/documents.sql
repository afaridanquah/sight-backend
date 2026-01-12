-- name: InsertDocument :exec
INSERT INTO documents (id, customer_id, business_id, original_name, filename, mimetype, created_at, status, creator_id, updated_at)
VALUES(@id, @customer_id, @business_id, @original_name, @filename, @mimetype, @created_at, @status, @creator_id, @updated_at);


-- name: UpdateDocument :exec
UPDATE documents
SET status = @status
WHERE id = @id;
