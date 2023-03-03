package models

// Hello represents message from one person to another
type Hello struct {
	ID      int64  `db:"id"`
	From    string `db:"sender"`
	To      string `db:"recipient"`
	Message string `db:"message"`
}
