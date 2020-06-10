package domain

import (
	"database/sql"
	"time"
)

//d3:entity
//d3_table:lw_wish
type Wish struct {
	id       sql.NullInt64 `d3:"pk:auto"`
	content  string
	createAt time.Time
}

func newWish(content string) *Wish {
	return &Wish{
		content:  content,
		createAt: time.Now(),
	}
}

func (w *Wish) CreateAt() time.Time {
	return w.createAt
}

func (w *Wish) Content() string {
	return w.content
}

func (w *Wish) grant() {

}
