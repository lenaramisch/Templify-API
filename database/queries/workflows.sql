-- name: AddWorkflow :one
INSERT INTO workflows (
    name, email_template_name, email_subject, static_attachments, templated_pdfs
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: GetWorkflowByName :one
SELECT * FROM workflows
WHERE name = $1 LIMIT 1;
