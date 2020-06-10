package main

import (
	"context"
	"github.com/godzie44/d3/adapter"
	"github.com/godzie44/d3/orm"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"log"
	"lw/internal/application"
	"lw/internal/domain"
	"lw/internal/infrastructure/repository"
	"net/http"
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

	d3rep, err := d3orm.MakeRepository(&domain.User{})
	if err != nil {
		log.Fatal(err.Error())
	}

	userService := application.NewUserService(repository.NewD3UserRepo(d3rep))
	newUser := &newUserHandler{service: userService}

	r := mux.NewRouter()
	r.Handle("/user", newUser).Methods("POST")
	r.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("pong"))
	}).Methods("POST")
	r.Use(makeOrmMiddleware(d3orm))

	log.Fatal(http.ListenAndServe("localhost:8089", r))
}

func makeOrmMiddleware(d3orm *orm.Orm) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := d3orm.CtxWithSession(context.Background())
			next.ServeHTTP(w, r.WithContext(ctx))
			if err := orm.Session(ctx).Flush(); err != nil {
				w.WriteHeader(500)
				_, _ = w.Write([]byte("internal error"))
			}
		})
	}
}

type newUserHandler struct {
	service *application.UserService
}

func (n *newUserHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	_ = request.ParseForm()
	err := n.service.NewUser(request.Context(), request.Form.Get("name"), request.Form.Get("email"))
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte(err.Error()))
	}
}
