-- name: AddPDFTemplate :one
INSERT INTO pdftemplates (
    name, templ_string
) VALUES (
    $1, $2
)
RETURNING *;

-- name: GetPDFTemplateByName :one
SELECT * FROM pdftemplates
WHERE name = $1 LIMIT 1;
