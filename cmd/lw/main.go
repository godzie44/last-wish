package main

import (
	"context"
	"encoding/json"
	"github.com/godzie44/d3/adapter"
	d3pgx "github.com/godzie44/d3/adapter/pgx"
	"github.com/godzie44/d3/orm"
	"github.com/godzie44/lw/internal/application"
	"github.com/godzie44/lw/internal/domain"
	appinfr "github.com/godzie44/lw/internal/infrastructure/application"
	"github.com/godzie44/lw/internal/infrastructure/domain/repository"
	"github.com/godzie44/lw/internal/infrastructure/domain/service"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
	"log"
	"net/http"
	"strconv"
)

func main() {
	pg, err := pgx.Connect(context.Background(), "postgres://postgres:postgres@0.0.0.0:5432/lw")
	if err != nil {
		log.Fatal(err.Error())
	}

	d3orm := orm.NewOrm(d3pgx.NewGoPgXAdapter(pg, &adapter.SquirrelAdapter{}))
	if err := d3orm.Register(&domain.User{}, &domain.Wish{}); err != nil {
		log.Fatal(err.Error())
	}

	d3rep, err := d3orm.MakeRepository(&domain.User{})
	if err != nil {
		log.Fatal(err.Error())
	}

	userService := appinfr.NewTransactionalUserService(
		application.NewUserService(repository.NewD3UserRepo(d3rep), &service.NotifyService{}),
	)
	wishService := appinfr.NewTransactionalWishService(
		application.NewWishService(repository.NewD3UserRepo(d3rep)),
	)

	r := mux.NewRouter()
	{
		r.HandleFunc("/user", func(writer http.ResponseWriter, request *http.Request) {
			_ = request.ParseForm()
			id, err := userService.NewUser(request.Context(), request.Form.Get("name"), request.Form.Get("email"))
			if err != nil {
				handleError(err, writer)
			} else {
				handleOk(struct {
					ID int64
				}{
					ID: id,
				}, writer)
			}
		}).Methods("POST")

		r.HandleFunc("/wish", func(writer http.ResponseWriter, request *http.Request) {
			_ = request.ParseForm()
			userId, _ := strconv.Atoi(request.Form.Get("userId"))
			err := wishService.NewWish(request.Context(), int64(userId), request.Form.Get("content"))
			if err != nil {
				handleError(err, writer)
			} else {
				handleOk(nil, writer)
			}
		}).Methods("POST")

		r.HandleFunc("/user/friend", func(writer http.ResponseWriter, request *http.Request) {
			_ = request.ParseForm()
			userId, _ := strconv.Atoi(request.Form.Get("userId"))
			friendId, _ := strconv.Atoi(request.Form.Get("friendId"))
			err := userService.AddFriend(request.Context(), int64(userId), int64(friendId))
			if err != nil {
				handleError(err, writer)
			} else {
				handleOk(nil, writer)
			}
		}).Methods("POST")

		r.HandleFunc("/user/release", func(writer http.ResponseWriter, request *http.Request) {
			_ = request.ParseForm()
			userId, _ := strconv.Atoi(request.Form.Get("userId"))
			err := userService.ReleaseWishes(request.Context(), int64(userId))
			if err != nil {
				handleError(err, writer)
			} else {
				handleOk(nil, writer)
			}
		}).Methods("POST")

		r.Use(makeOrmMiddleware(d3orm))
	}

	log.Fatal(http.ListenAndServe("localhost:8089", r))
}

func makeOrmMiddleware(d3orm *orm.Orm) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := d3orm.CtxWithSession(context.Background())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func handleError(err error, writer http.ResponseWriter) {
	type response struct {
		Error string `json:"error"`
	}

	var resp = response{}
	writer.WriteHeader(500)
	resp.Error = err.Error()

	jsonResp, _ := json.Marshal(resp)
	writer.Write(jsonResp)
}

func handleOk(data interface{}, writer http.ResponseWriter) {
	type response struct {
		Result interface{} `json:"result"`
	}

	jsonResp, _ := json.Marshal(response{Result: data})
	writer.Write(jsonResp)
}
