-- name: AddSMSTemplate :one
INSERT INTO smstemplates (
    name, templ_string
) VALUES (
    $1, $2
)
RETURNING *;

-- name: GetSMSTemplateByName :one
SELECT * FROM smstemplates
WHERE name = $1 LIMIT 1;
