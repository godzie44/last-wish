package domain

import (
	"database/sql"
	"errors"
	"log"
	"time"
	"unicode/utf8"
)

//d3:entity
//d3_table:lw_wish
type Wish struct {
	id       sql.NullInt64 `d3:"pk:auto"`
	content  string
	createAt time.Time
}

var ErrTooBigWish = errors.New("wish content is too big")

func newWish(content string) (*Wish, error) {
	if utf8.RuneCountInString(content) > 255 {
		return nil, ErrTooBigWish
	}

	return &Wish{
		content:  content,
		createAt: time.Now(),
	}, nil
}

func (w *Wish) ID() int64 {
	return w.id.Int64
}

func (w *Wish) CreateAt() time.Time {
	return w.createAt
}

func (w *Wish) Content() string {
	return w.content
}

func (w *Wish) grant() {
	log.Printf("grant wish: %s", w.content)
}
