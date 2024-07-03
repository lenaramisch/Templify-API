package db

type Template struct {
	ID         int    `db:"id"`
	Name       string `db:"name"`
	MJMLString string `db:"mjml_string"`
	CreatedAt  string `db:"created_at"`
}
