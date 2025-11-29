-- name: GetCatalogItems :many
SELECT id, title, description, disabled
FROM catalog
ORDER BY title;
