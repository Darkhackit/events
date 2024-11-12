package api

import (
	"context"
	db "github.com/Darkhackit/events/db/sqlc"
	"github.com/Darkhackit/events/repository"
	"github.com/Darkhackit/events/service"
	"github.com/Darkhackit/events/token"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
)

func Start() {
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, "postgresql://root:password@localhost:5432/test_event?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	queries := db.New(conn)

	PasetoToken, err := token.NewPasetoToken()
	if err != nil {
		panic(err)
	}

	uh := UserHandler{service: service.NewUserService(repository.NewUserRepositoryDB(queries, PasetoToken))}

	router := mux.NewRouter()

	router.HandleFunc("/users", uh.CreateUser).Methods("POST")
	router.HandleFunc("/login", uh.LoginUser).Methods("POST")

	protectedRouter := router.PathPrefix("/").Subrouter()
	protectedRouter.Use(AuthMiddleware(PasetoToken))
	protectedRouter.HandleFunc("/users", uh.GetUsers).Methods("GET")

	err = http.ListenAndServe(":8000", router)
	if err != nil {
		panic(err)
	}
}
