package api_test

import (
	"avito-winter-2025/internal/delivery"
	"avito-winter-2025/internal/entity"
	"avito-winter-2025/internal/repo"
	"avito-winter-2025/internal/usecase"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

var userKey string = "user"

func InitPostgres(user, password, host, port, name string) *pgxpool.Pool {
	PG_CONN := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, name)
	db, err := pgxpool.New(context.Background(), PG_CONN)
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL", err)
		return nil
	}
	if err := db.Ping(context.Background()); err != nil {
		log.Fatal("Failed to connect to PostgreSQL", err)
		return nil
	}
	return db
}

type ShopTestSuite struct {
	db      *pgxpool.Pool
	handler *delivery.ShopHandler
	url     string
	user    entity.User
	suite.Suite
}

func (s *ShopTestSuite) SetupTest() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(err)
	}
	db := InitPostgres(os.Getenv("TEST_DB_USER"), os.Getenv("TEST_DB_PASSWORD"), os.Getenv("TEST_DB_HOST"), os.Getenv("TEST_DB_PORT"), os.Getenv("TEST_DB_NAME"))
	if db == nil {
		s.T().Fatal("Failed to initialize database connection")
	}
	s.db = db
	merchRepo := repo.NewMerch(db)
	userRepo := repo.NewUser(db)
	coinRepo := repo.NewCoin(db)
	merchUC := usecase.NewMerch(merchRepo, coinRepo)
	userUC := usecase.NewUser(userRepo)
	coinUC := usecase.NewCoin(coinRepo, userRepo)
	s.handler = delivery.NewShopHandler(merchUC, userUC, coinUC)
	s.url = "/buy"
	s.user = entity.User{
		Name:  "sofia",
		Coins: 1000,
	}
}

func (s *ShopTestSuite) TestBuyMerch_Success() {
	query := `INSERT INTO "user" (name, password, coins) VALUES ($1, $2, $3)
				RETURNING id, name, coins`
	ctx := context.Background()
	s.db.Exec(ctx, `DELETE FROM "user"`)
	s.db.QueryRow(ctx, query, s.user.Name, "12345", s.user.Coins).Scan(&s.user.ID, &s.user.Name, &s.user.Coins)

	merch := "t-shirt"
	req := httptest.NewRequest("GET", s.url+"/"+merch, nil)
	req = req.WithContext(context.WithValue(req.Context(), userKey, s.user))

	rw := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc(s.url+"/{item}", s.handler.BuyMerch)
	router.ServeHTTP(rw, req)

	s.Equal(http.StatusOK, rw.Code)

	var newBalance uint32
	var cost uint32
	s.db.QueryRow(context.Background(), `SELECT coins FROM "user" WHERE id = $1`, s.user.ID).Scan(&newBalance)
	s.db.QueryRow(context.Background(), `SELECT cost FROM merch WHERE name = $1`, merch).Scan(&cost)
	s.Equal(s.user.Coins-cost, newBalance)
}

func (s *ShopTestSuite) TestBuyMerch_NoAuth() {
	merch := "t-shirt"
	req := httptest.NewRequest("GET", s.url+"/"+merch, nil)
	req = req.WithContext(context.WithValue(req.Context(), "", s.user))

	rw := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc(s.url+"/{item}", s.handler.BuyMerch)
	router.ServeHTTP(rw, req)

	s.Equal(http.StatusUnauthorized, rw.Code)
}

func (s *ShopTestSuite) TestBuyMerch_WrongMerch() {
	query := `INSERT INTO "user" (name, password, coins) VALUES ($1, $2, $3)
				RETURNING id, name, coins`
	ctx := context.Background()
	s.db.Exec(ctx, `DELETE FROM "user"`)
	s.db.QueryRow(ctx, query, s.user.Name, "12345", s.user.Coins).Scan(&s.user.ID, &s.user.Name, &s.user.Coins)

	merch := "t-shi"
	req := httptest.NewRequest("GET", s.url+"/"+merch, nil)
	req = req.WithContext(context.WithValue(req.Context(), userKey, s.user))

	rw := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc(s.url+"/{item}", s.handler.BuyMerch)
	router.ServeHTTP(rw, req)

	s.Equal(http.StatusBadRequest, rw.Code)
	bodyBytes, err := io.ReadAll(rw.Body)
	s.Equal(err, nil)
	body := string(bodyBytes)
	s.Equal(body, `{"error":"Мерч не найден"}`)
}

func (s *ShopTestSuite) TestBuyMerch_NotEnoughCoin() {
	merch := "t-shirt"
	query := `INSERT INTO "user" (name, password, coins) VALUES ($1, $2, $3)
				RETURNING id, name, coins`
	var cost uint32
	ctx := context.Background()
	s.db.QueryRow(ctx, `SELECT cost FROM merch WHERE name = $1`, merch).Scan(&cost)
	s.db.Exec(ctx, `DELETE FROM "user"`)
	s.db.QueryRow(ctx, query, s.user.Name, "12345", cost-1).Scan(&s.user.ID, &s.user.Name, &s.user.Coins)

	req := httptest.NewRequest("GET", s.url+"/"+merch, nil)
	req = req.WithContext(context.WithValue(req.Context(), userKey, s.user))

	rw := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc(s.url+"/{item}", s.handler.BuyMerch)
	router.ServeHTTP(rw, req)

	s.Equal(http.StatusBadRequest, rw.Code)
	bodyBytes, err := io.ReadAll(rw.Body)
	s.Equal(err, nil)
	body := string(bodyBytes)
	s.Equal(body, `{"error":"У вас недостаточно стредств"}`)
}

func (s *ShopTestSuite) TestBuyMerch_ServerError() {
	query := `INSERT INTO "user" (name, password, coins) VALUES ($1, $2, $3)
				RETURNING id, name, coins`
	ctx := context.Background()
	s.db.Exec(ctx, `DELETE FROM "user"`)
	s.db.QueryRow(ctx, query, s.user.Name, "12345", s.user.Coins).Scan(&s.user.ID, &s.user.Name, &s.user.Coins)

	merch := "t-shirt"
	req := httptest.NewRequest("GET", s.url+"/"+merch, nil)
	s.user.ID = 0
	req = req.WithContext(context.WithValue(req.Context(), userKey, s.user))

	rw := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc(s.url+"/{item}", s.handler.BuyMerch)
	router.ServeHTTP(rw, req)

	s.Equal(http.StatusInternalServerError, rw.Code)
	bodyBytes, err := io.ReadAll(rw.Body)
	s.Equal(err, nil)
	body := string(bodyBytes)
	s.Equal(body, `{"error":"Ошибка сервера"}`)
}

func TestShopTestSuite(t *testing.T) {
	suite.Run(t, new(ShopTestSuite))
}
