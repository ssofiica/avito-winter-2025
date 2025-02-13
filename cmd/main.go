package main

import (
	"avito-winter-2025/config"
	"avito-winter-2025/internal/delivery"
	"avito-winter-2025/internal/repo"
	"avito-winter-2025/internal/usecase"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(err)
	}
}

func main() {
	cfg := config.Load()
	PG_CONN := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", os.Getenv("DATABASE_USER"), os.Getenv("DATABASE_PASSWORD"), os.Getenv("DATABASE_HOST"), os.Getenv("DATABASE_PORT"), os.Getenv("DATABASE_NAME"))
	db, err := pgxpool.New(context.Background(), PG_CONN)
	if err != nil {
		fmt.Println("error wih db", err)
	}
	defer db.Close()

	userRepo := repo.NewUser(db)
	coinRepo := repo.NewCoin(db)
	merchRepo := repo.NewMerch(db)

	userUsecase := usecase.NewUser(userRepo)
	coinUsecase := usecase.NewCoin(coinRepo)
	merchUsecase := usecase.NewMerch(merchRepo, coinRepo)

	authHandler := delivery.NewAuthHandler(userUsecase)
	coinHandler := delivery.NewCoinHandler(coinUsecase, userUsecase)
	shopHandler := delivery.NewShopHandler(merchUsecase)

	r := mux.NewRouter().PathPrefix("/api").Subrouter()

	r.HandleFunc("/ok", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("query to path: " + r.URL.String())
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	//r.Use() //тут мидлвары вставить надо аргументами
	// r.HandleFunc("/info").Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/sendCoin", delivery.JWTMiddleware(coinHandler.SendCoin)).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/buy/{item}", delivery.JWTMiddleware(shopHandler.BuyMerch)).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/auth", authHandler.Auth).Methods(http.MethodPost, http.MethodOptions)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")),
		Handler:           r,
		ReadTimeout:       cfg.Server.ReadTimeout,
		WriteTimeout:      cfg.Server.WriteTimeout,
		ReadHeaderTimeout: cfg.Server.ReadHeaderTimeout,
		IdleTimeout:       cfg.Server.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("listen: %s\\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown err:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}
