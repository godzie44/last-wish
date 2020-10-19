package main

import (
	"context"
	"encoding/json"
	"github.com/godzie44/d3/adapter/pgx"
	"github.com/godzie44/d3/orm"
	"github.com/godzie44/lw/internal/application"
	"github.com/godzie44/lw/internal/domain"
	appinfr "github.com/godzie44/lw/internal/infrastructure/application"
	"github.com/godzie44/lw/internal/infrastructure/domain/repository"
	"github.com/godzie44/lw/internal/infrastructure/domain/service"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
)

func main() {
	cfg, err := pgxpool.ParseConfig(os.Getenv("PG_DSN"))
	if err != nil {
		log.Fatal(err.Error())
	}
	driver, err := pgx.NewPgxPoolDriver(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer driver.Close()

	d3orm := orm.New(driver)
	if err := d3orm.Register(&domain.User{}, &domain.Wish{}); err != nil {
		log.Fatal(err.Error())
	}

	if err = createSchema(d3orm, driver.UnwrapConn().(*pgxpool.Pool)); err != nil {
		log.Fatal(err.Error())
	}

	var userRepository domain.UserRepository
	{
		d3UserRep, err := d3orm.MakeRepository(&domain.User{})
		if err != nil {
			log.Fatal(err.Error())
		}
		userRepository = repository.NewUserRepository(d3UserRep)
	}

	var wishRepository domain.WishRepository
	{
		d3WishRep, err := d3orm.MakeRepository(&domain.Wish{})
		if err != nil {
			log.Fatal(err.Error())
		}
		wishRepository = repository.NewWishRepository(d3WishRep)
	}

	userService := appinfr.NewTransactionalUserService(
		application.NewUserService(userRepository, &service.NotifyService{}),
	)
	wishService := appinfr.NewTransactionalWishService(
		application.NewWishService(userRepository, wishRepository),
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

		r.HandleFunc("/user/{id}/wish", func(writer http.ResponseWriter, request *http.Request) {
			vars := mux.Vars(request)
			_ = request.ParseForm()
			userId, _ := strconv.Atoi(vars["id"])
			err := wishService.NewWish(request.Context(), int64(userId), request.Form.Get("content"))
			if err != nil {
				handleError(err, writer)
			} else {
				handleOk(nil, writer)
			}
		}).Methods("POST")

		r.HandleFunc("/user/{user_id}/friend/{friend_id}", func(writer http.ResponseWriter, request *http.Request) {
			vars := mux.Vars(request)
			userId, _ := strconv.Atoi(vars["user_id"])
			friendId, _ := strconv.Atoi(vars["friend_id"])
			err := userService.AddFriend(request.Context(), int64(userId), int64(friendId))
			if err != nil {
				handleError(err, writer)
			} else {
				handleOk(nil, writer)
			}
		}).Methods("POST")

		r.HandleFunc("/user/{user_id}/wish/release", func(writer http.ResponseWriter, request *http.Request) {
			vars := mux.Vars(request)
			userId, _ := strconv.Atoi(vars["user_id"])
			err := userService.ReleaseWishes(request.Context(), int64(userId))
			if err != nil {
				handleError(err, writer)
			} else {
				handleOk(nil, writer)
			}
		}).Methods("POST")

		r.HandleFunc("/user/{user_id}/wish/{wish_id}", func(writer http.ResponseWriter, request *http.Request) {
			vars := mux.Vars(request)
			userId, _ := strconv.Atoi(vars["user_id"])
			wishId, _ := strconv.Atoi(vars["wish_id"])
			err := wishService.DeleteWish(request.Context(), int64(userId), int64(wishId))
			if err != nil {
				handleError(err, writer)
			} else {
				handleOk(nil, writer)
			}
		}).Methods("DELETE")

		r.HandleFunc("/user/{user_id}/wish/{wish_id}", func(writer http.ResponseWriter, request *http.Request) {
			vars := mux.Vars(request)
			userId, _ := strconv.Atoi(vars["user_id"])
			wishId, _ := strconv.Atoi(vars["wish_id"])
			_ = request.ParseForm()
			err := wishService.UpdateWish(request.Context(), int64(userId), int64(wishId), request.Form.Get("content"))
			if err != nil {
				handleError(err, writer)
			} else {
				handleOk(nil, writer)
			}
		}).Methods("PUT")

		r.Use(makeOrmMiddleware(d3orm))
	}

	log.Fatal(http.ListenAndServe("0.0.0.0:8080", r))
}

func makeOrmMiddleware(d3orm *orm.Orm) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := d3orm.CtxWithSession(r.Context())
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

func createSchema(orm *orm.Orm, conn *pgxpool.Pool) error {
	sql, err := orm.GenerateSchema()
	if err != nil {
		return err
	}

	_, err = conn.Exec(context.Background(), sql)
	return err
}
