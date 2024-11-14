package api

import (
	"context"
	"errors"
	db "github.com/Darkhackit/events/db/sqlc"
	"github.com/Darkhackit/events/repository"
	"github.com/Darkhackit/events/service"
	"github.com/Darkhackit/events/token"
	"github.com/Darkhackit/events/worker"
	"github.com/gorilla/mux"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
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
	redisOpt := asynq.RedisClientOpt{
		Addr: "localhost:6379",
	}
	distributor := worker.NewRedisTaskDistributor(redisOpt)
	uh := UserHandler{service: service.NewUserService(repository.NewUserRepositoryDB(queries, PasetoToken, distributor))}

	router := mux.NewRouter()

	router.HandleFunc("/users", uh.CreateUser).Methods("POST")
	router.HandleFunc("/login", uh.LoginUser).Methods("POST")

	protectedRouter := router.PathPrefix("/").Subrouter()
	protectedRouter.Use(AuthMiddleware(PasetoToken))
	protectedRouter.HandleFunc("/users", uh.GetUsers).Methods("GET")

	//go RunTaskProcessor(redisOpt, queries)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, *queries)
	go func() {
		log.Info().Msg("Starting task processor")
		if err := taskProcessor.Start(); err != nil {
			log.Fatal().Msg(err.Error())
		}
	}()
	srv := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		log.Info().Msg("Shutting down server...")
		if err := srv.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("Error during shutdown")
		}
		taskProcessor.Stop() // Clean up task processor
	}()
	log.Info().Msg("Server is running on :8000")
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().Err(err).Msg("ListenAndServe error")
	}
}

func RunTaskProcessor(redisOpt asynq.RedisClientOpt, store *db.Queries) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, *store)
	log.Info().Msg("Starting task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
}
