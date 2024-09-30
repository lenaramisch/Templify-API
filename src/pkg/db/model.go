package db

type Workflow struct {
	ID                int    `db:"id"`
	Name              string `db:"name"`
	EmailTemplateName string `db:"email_template_name"`
	EmailSubject      string `db:"email_subject"`
	StaticAttachments string `db:"static_attachments"`
	TemplatedPDFs     string `db:"templated_pdfs"`
	CreatedAt         string `db:"created_at"`
}

type Template struct {
	ID          int    `db:"id"`
	Name        string `db:"name"`
	TemplString string `db:"templ_string"`
	IsMJML      bool   `db:"is_mjml"`
	CreatedAt   string `db:"created_at"`
}
