-- name: SaveData :exec
INSERT INTO data (key, value)
VALUES ($1, $2);
