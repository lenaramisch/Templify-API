package db

type Template struct {
	ID          int    `db:"id"`
	Name        string `db:"name"`
	TemplString string `db:"templ_string"`
	IsMJML      bool   `db:"is_mjml"`
	CreatedAt   string `db:"created_at"`
}

type PDFTemplate struct {
	ID          int    `db:"id"`
	Name        string `db:"name"`
	TypstString string `db:"typst_string"`
	CreatedAt   string `db:"created_at"`
}

type SMSTemplate struct {
	ID          int    `db:"id"`
	Name        string `db:"name"`
	TemplString string `db:"templ_string"`
	CreatedAt   string `db:"created_at"`
}
