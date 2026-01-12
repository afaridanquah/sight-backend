-- name: InsertOrg :exec
INSERT INTO organizations (id,name,status,created_at,updated_at)
VALUES (@id,@name,@status,@created_at,@updated_at);
