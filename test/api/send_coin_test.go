package api_test

import (
	"avito-winter-2025/internal/delivery"
	"avito-winter-2025/internal/entity"
	"avito-winter-2025/internal/repo"
	"avito-winter-2025/internal/usecase"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

type SendCoinTestSuite struct {
	db       *pgxpool.Pool
	handler  *delivery.CoinHandler
	url      string
	fromUser entity.User
	suite.Suite
}

func (s *SendCoinTestSuite) SetupTest() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(err)
	}
	db := InitPostgres(os.Getenv("TEST_DB_USER"), os.Getenv("TEST_DB_PASSWORD"), os.Getenv("TEST_DB_HOST"), os.Getenv("TEST_DB_PORT"), os.Getenv("TEST_DB_NAME"))
	if db == nil {
		s.T().Fatal("Failed to initialize database connection")
	}
	s.db = db
	userRepo := repo.NewUser(db)
	coinRepo := repo.NewCoin(db)
	userUC := usecase.NewUser(userRepo)
	coinUC := usecase.NewCoin(coinRepo, userRepo)
	s.handler = delivery.NewCoinHandler(coinUC, userUC)
	s.url = "/sendCoin"
	s.fromUser = entity.User{
		Name:  "sofia",
		Coins: 1000,
	}
}

func (s *SendCoinTestSuite) TestSendCoin_Success() {
	toUser := entity.User{
		Name:  "mary",
		Coins: 1000,
	}
	amount := 200
	query := `INSERT INTO "user" (name, password, coins) VALUES ($1, $2, $3)
				RETURNING id, name, coins`
	ctx := context.Background()
	s.db.Exec(ctx, `DELETE FROM "user"`)
	s.db.QueryRow(ctx, query, s.fromUser.Name, "12345", s.fromUser.Coins).
		Scan(&s.fromUser.ID, &s.fromUser.Name, &s.fromUser.Coins)
	s.db.QueryRow(ctx, query, toUser.Name, "12345", toUser.Coins).
		Scan(&toUser.ID, &toUser.Name, &toUser.Coins)

	data := []byte(`{"toUser":"` + toUser.Name + `", "amount":` + strconv.Itoa(amount) + `}`)
	req := httptest.NewRequest("POST", s.url, bytes.NewBuffer(data))
	req = req.WithContext(context.WithValue(req.Context(), userKey, s.fromUser))

	rw := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc(s.url, s.handler.SendCoin)
	router.ServeHTTP(rw, req)

	s.Equal(http.StatusOK, rw.Code)

	var fromUserBalance uint32
	var toUserBalance uint32
	query1 := `SELECT coins FROM "user" WHERE id = $1`
	s.db.QueryRow(context.Background(), query1, s.fromUser.ID).Scan(&fromUserBalance)
	s.db.QueryRow(context.Background(), query1, toUser.ID).Scan(&toUserBalance)
	s.Equal(fromUserBalance+uint32(amount), s.fromUser.Coins)
	s.Equal(toUserBalance-uint32(amount), toUser.Coins)
}

func (s *SendCoinTestSuite) TestSendCoin_NoAuth() {
	req := httptest.NewRequest("POST", s.url, nil)
	req = req.WithContext(context.WithValue(req.Context(), "", s.fromUser))

	rw := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc(s.url, s.handler.SendCoin)
	router.ServeHTTP(rw, req)

	s.Equal(http.StatusUnauthorized, rw.Code)
}

func (s *SendCoinTestSuite) TestSendCoin_NoRequestData() {
	req := httptest.NewRequest("POST", s.url, nil)
	req = req.WithContext(context.WithValue(req.Context(), userKey, s.fromUser))

	rw := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc(s.url, s.handler.SendCoin)
	router.ServeHTTP(rw, req)

	s.Equal(http.StatusBadRequest, rw.Code)
}

func (s *SendCoinTestSuite) TestSendCoin_InvalidRequestData() {
	data := []byte(`{"toUser":" ", "amount":-1}`)
	req := httptest.NewRequest("POST", s.url, bytes.NewBuffer(data))
	req = req.WithContext(context.WithValue(req.Context(), userKey, s.fromUser))

	rw := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc(s.url, s.handler.SendCoin)
	router.ServeHTTP(rw, req)

	s.Equal(http.StatusBadRequest, rw.Code)
}

func (s *SendCoinTestSuite) TestSendCoin_NoDestinationUser() {
	toUser := entity.User{
		Name: "mary",
	}
	amount := 200
	query := `INSERT INTO "user" (name, password, coins) VALUES ($1, $2, $3)
				RETURNING id, name, coins`
	ctx := context.Background()
	s.db.Exec(ctx, `DELETE FROM "user"`)
	s.db.QueryRow(ctx, query, s.fromUser.Name, "12345", s.fromUser.Coins).
		Scan(&s.fromUser.ID, &s.fromUser.Name, &s.fromUser.Coins)

	data := []byte(`{"toUser":"` + toUser.Name + `", "amount":` + strconv.Itoa(amount) + `}`)
	req := httptest.NewRequest("POST", s.url, bytes.NewBuffer(data))
	req = req.WithContext(context.WithValue(req.Context(), userKey, s.fromUser))

	rw := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc(s.url, s.handler.SendCoin)
	router.ServeHTTP(rw, req)

	s.Equal(http.StatusBadRequest, rw.Code)
	bodyBytes, err := io.ReadAll(rw.Body)
	s.Equal(err, nil)
	body := string(bodyBytes)
	s.Equal(body, `{"error":"Пользователь не найден"}`)
}

func (s *SendCoinTestSuite) TestSendCoin_NotEnoughCoins() {
	toUser := entity.User{
		Name:  "mary",
		Coins: 1000,
	}
	amount := 200
	query := `INSERT INTO "user" (name, password, coins) VALUES ($1, $2, $3)
				RETURNING id, name, coins`
	ctx := context.Background()
	s.db.Exec(ctx, `DELETE FROM "user"`)
	s.db.QueryRow(ctx, query, s.fromUser.Name, "12345", amount-1).
		Scan(&s.fromUser.ID, &s.fromUser.Name, &s.fromUser.Coins)
	s.db.QueryRow(ctx, query, toUser.Name, "12345", toUser.Coins).
		Scan(&toUser.ID, &toUser.Name, &toUser.Coins)

	data := []byte(`{"toUser":"` + toUser.Name + `", "amount":` + strconv.Itoa(amount) + `}`)
	req := httptest.NewRequest("POST", s.url, bytes.NewBuffer(data))
	req = req.WithContext(context.WithValue(req.Context(), userKey, s.fromUser))

	rw := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc(s.url, s.handler.SendCoin)
	router.ServeHTTP(rw, req)

	s.Equal(http.StatusBadRequest, rw.Code)
	bodyBytes, err := io.ReadAll(rw.Body)
	s.Equal(err, nil)
	body := string(bodyBytes)
	s.Equal(body, `{"error":"У вас недостаточно стредств"}`)
}

func TestSendCoinTestSuite(t *testing.T) {
	suite.Run(t, new(SendCoinTestSuite))
}
