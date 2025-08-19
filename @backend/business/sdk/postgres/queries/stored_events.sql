-- name: InsertEvent :exec
    INSERT INTO stored_events(id, type, aggregate_id, aggregate_type, aggregate_version, data, meta_data, occured_at, registered_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);
