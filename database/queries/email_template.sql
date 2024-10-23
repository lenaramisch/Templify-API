-- name: GetEmailTemplateByName :one
SELECT * FROM emailtemplates
WHERE name = $1 LIMIT 1;

-- name: AddEmailTemplate :one
INSERT INTO emailtemplates (
    name, templ_string, is_mjml
) VALUES (
    $1, $2, $3
)
RETURNING *;
