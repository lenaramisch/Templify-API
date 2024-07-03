package db

type Template struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	MJMLString string `json:"mjml_string"`
	CreatedAt  string `json:"created_at"`
}
