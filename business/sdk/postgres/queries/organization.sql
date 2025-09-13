-- name: InsertOrg :exec
INSERT INTO organizations (id,name,status)
VALUES ($1,$2,$3);
