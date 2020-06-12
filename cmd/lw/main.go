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
	"strconv"
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
	wishService := application.NewWishService(repository.NewD3UserRepo(d3rep))

	r := mux.NewRouter()
	r.Handle("/user", &newUserHandler{userService}).Methods("POST")
	r.Handle("/wish", &newWishHandler{wishService}).Methods("POST")
	r.Handle("/user/friend", &addFriendHandler{userService}).Methods("POST")
	r.Handle("/user/release", &releaseWishesHandler{userService}).Methods("POST")
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
	handleError(err, writer)
}

type newWishHandler struct {
	service *application.WishService
}

func (n *newWishHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	_ = request.ParseForm()
	userId, _ := strconv.Atoi(request.Form.Get("userId"))
	err := n.service.NewWish(request.Context(), int64(userId), request.Form.Get("content"))
	handleError(err, writer)
}

type addFriendHandler struct {
	service *application.UserService
}

func (n *addFriendHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	_ = request.ParseForm()
	userId, _ := strconv.Atoi(request.Form.Get("userId"))
	friendId, _ := strconv.Atoi(request.Form.Get("friendId"))
	err := n.service.AddFriend(request.Context(), int64(userId), int64(friendId))
	handleError(err, writer)
}

type releaseWishesHandler struct {
	service *application.UserService
}

func (n *releaseWishesHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	_ = request.ParseForm()
	userId, _ := strconv.Atoi(request.Form.Get("userId"))
	err := n.service.ReleaseWishes(request.Context(), int64(userId))
	handleError(err, writer)
}

func handleError(err error, writer http.ResponseWriter) {
	if err != nil {
		writer.WriteHeader(500)
		writer.Write([]byte(err.Error()))
	}
}
