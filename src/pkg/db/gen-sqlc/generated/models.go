// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package repo_sqlc

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Emailtemplate struct {
	ID          int32
	Name        pgtype.Text
	TemplString pgtype.Text
	IsMjml      pgtype.Bool
	CreatedAt   pgtype.Timestamptz
}

type Pdftemplate struct {
	ID          int32
	Name        pgtype.Text
	TemplString pgtype.Text
	CreatedAt   pgtype.Timestamptz
}

type Smstemplate struct {
	ID          int32
	Name        pgtype.Text
	TemplString pgtype.Text
	CreatedAt   pgtype.Timestamptz
}

type Workflow struct {
	ID                int32
	Name              pgtype.Text
	EmailTemplateName pgtype.Text
	EmailSubject      pgtype.Text
	StaticAttachments pgtype.Text
	TemplatedPdfs     pgtype.Text
	CreatedAt         pgtype.Timestamptz
}
