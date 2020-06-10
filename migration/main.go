package main

import (
	"context"
	"github.com/godzie44/d3/adapter"
	"github.com/godzie44/d3/orm"
	"github.com/jackc/pgx/v4"
	"log"
	"lw/internal/domain"
)

func main() {
	pg, err := pgx.Connect(context.Background(), "postgres://postgres:postgres@0.0.0.0:5432/lw")
	if err != nil {
		log.Fatal(err.Error())
	}

	d3orm := orm.NewOrm(adapter.NewGoPgXAdapter(pg, &adapter.SquirrelAdapter{}))
	if err := d3orm.Register(&domain.User{}, &domain.Wish{}); err != nil {
		log.Fatal(err.Error())
	}

	sql, err := d3orm.GenerateSchema()
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = pg.Exec(context.Background(), sql)
	if err != nil {
		log.Fatal(err.Error())
	}

}
