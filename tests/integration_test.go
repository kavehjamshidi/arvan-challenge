package tests

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/kavehjamshidi/arvan-challenge/api"
	"github.com/kavehjamshidi/arvan-challenge/api/handler"
	"github.com/kavehjamshidi/arvan-challenge/driver/db/postgres"
	redis2 "github.com/kavehjamshidi/arvan-challenge/driver/db/redis"
	"github.com/kavehjamshidi/arvan-challenge/external/rate_limiter"
	"github.com/kavehjamshidi/arvan-challenge/queue"
	"github.com/kavehjamshidi/arvan-challenge/repository/file"
	"github.com/kavehjamshidi/arvan-challenge/repository/user"
	"github.com/kavehjamshidi/arvan-challenge/service/rate_limit"
	"github.com/kavehjamshidi/arvan-challenge/service/upload"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const timeoutMiliseconds = 5000

type integrationTestSuite struct {
	suite.Suite
	redisClient *redis.Client
	db          *sql.DB
	server      *fiber.App
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, &integrationTestSuite{})
}

// Suite Setup
func (i *integrationTestSuite) SetupSuite() {
	viper.AutomaticEnv()
	viper.ReadInConfig()

	viper.SetDefault("TEST_SERVER_ADDRESS", ":4001")
	viper.SetDefault("TEST_DB_URI", "postgres://postgres:very-secret@127.0.0.1:5432/arvan?sslmode=disable")
	viper.SetDefault("TEST_REDIS_ADDRESS", "127.0.0.1:6379")
	viper.SetDefault("DB_MIGRATION_TABLE", "migrations")

	pg := postgres.TestSetup()
	postgres.TestMigrate(pg)
	i.db = pg

	redisClient := redis2.TestSetup()
	i.redisClient = redisClient

	rateLimiter := rate_limiter.NewRedisRateLimiter(redisClient)
	userRepo := user.NewPostgresUserRepository(pg)
	fileRepo := file.NewRedisFileRepository(redisClient)
	noOpQueue := queue.NewNoOpQueue()

	rateLimitSVC := rate_limit.NewRateLimitService(userRepo, rateLimiter)
	uploadSVC := upload.NewUploadService(fileRepo, userRepo, noOpQueue)

	rateLimitHandler := handler.NewRateLimitHandler(rateLimitSVC)
	uploadHandler := handler.NewUploadHandler(uploadSVC)

	server := api.SetupServer(rateLimitHandler, uploadHandler)

	i.server = server
}

func (i *integrationTestSuite) TearDownSuite() {
	i.db.Close()
	i.redisClient.Close()
	i.server.Shutdown()
}

func (i *integrationTestSuite) TearDownTest() {
	i.db.Exec("DELETE FROM user_usage WHERE user_id = '123456';")
	i.db.Exec("DELETE FROM users WHERE id = '123456';")

	i.redisClient.Del(context.TODO(), "rate:123456")
	i.redisClient.Del(context.TODO(), "quota:123456")
	i.redisClient.Del(context.TODO(), "file_id:abcdef")
	i.redisClient.Del(context.TODO(), "file_id:duplicate")
}

func (i *integrationTestSuite) BeforeTest(suiteName, testName string) {
	_, err := i.db.Exec(`INSERT INTO users(id, rate_limit, quota, created_at, updated_at)
VALUES('123456', 2, 10, NOW(), NOW()) ON CONFLICT (id) DO UPDATE SET rate_limit = 2;`)
	i.NoError(err)

	_, err = i.db.Exec(`INSERT INTO user_usage(user_id, quota, quota_usage,
                       start_date, end_date, created_at, updated_at) 
VALUES('123456', 10, 0, NOW(), NOW() + interval '1 month', NOW(), NOW()) ON CONFLICT (user_id) DO UPDATE SET quota_usage = 0;`)
	i.NoError(err)

	switch testName {
	case "TestRateLimiterFail":
		now := time.Now().UnixNano()
		err = i.redisClient.ZAdd(context.TODO(), "rate:123456", redis.Z{
			Score:  float64(now),
			Member: now,
		}).Err()
		i.NoError(err)

		err = i.redisClient.ZAdd(context.TODO(), "rate:123456", redis.Z{
			Score:  float64(now + 1),
			Member: now + 1,
		}).Err()
		i.NoError(err)

		err = i.redisClient.ZAdd(context.TODO(), "rate:123456", redis.Z{
			Score:  float64(now + 2),
			Member: now + 2,
		}).Err()
		i.NoError(err)
	case "TestUploadFail", "TestUploadSuccess":
		err = i.redisClient.Set(context.TODO(), "file_id:duplicate", 1, 0).Err()
		i.NoError(err)
	}
}

// Tests
func (i *integrationTestSuite) TestUploadValidationFail() {
	i.Run("fail - no file", func() {
		req := i.newHTTPRequest("123456", "", nil)

		res, err := i.server.Test(req)
		i.NoError(err)
		i.Equal(http.StatusBadRequest, res.StatusCode)

		body, err := io.ReadAll(res.Body)
		i.NoError(err)

		data := make(map[string]any)
		err = json.Unmarshal(body, &data)
		i.NoError(err)

		msg, _ := data["message"].(string)
		errMsg, _ := data["error"].(string)

		i.Equal("validation error", msg)
		i.Contains(errMsg, "could not find file in form data")
	})

	i.Run("fail - no fileID", func() {
		req := i.newHTTPRequest("123456", "", []byte("fake data"))

		res, err := i.server.Test(req)
		i.NoError(err)
		i.Equal(http.StatusBadRequest, res.StatusCode)

		body, err := io.ReadAll(res.Body)
		i.NoError(err)

		data := make(map[string]any)
		err = json.Unmarshal(body, &data)
		i.NoError(err)

		msg, _ := data["message"].(string)
		errMsg, _ := data["error"].(string)

		i.Equal("validation error", msg)
		i.Contains(errMsg, "file_id is required")
	})
}

func (i *integrationTestSuite) TestUploadFail() {
	i.Run("fail - duplicate fileID", func() {
		req := i.newHTTPRequest("123456", "duplicate", []byte("data"))

		res, err := i.server.Test(req)
		i.NoError(err)
		i.Equal(http.StatusConflict, res.StatusCode)

		body, err := io.ReadAll(res.Body)
		i.NoError(err)

		data := make(map[string]any)
		err = json.Unmarshal(body, &data)
		i.NoError(err)

		msg, _ := data["message"].(string)
		errMsg, _ := data["error"].(string)

		i.Equal("failed", msg)
		i.Contains(errMsg, "file_id already exists")
	})

	i.Run("fail - quota exceeded", func() {
		req := i.newHTTPRequest("123456", "abcdef", []byte("very big fake data"))

		res, err := i.server.Test(req)
		i.NoError(err)
		i.Equal(http.StatusForbidden, res.StatusCode)

		body, err := io.ReadAll(res.Body)
		i.NoError(err)

		data := make(map[string]any)
		err = json.Unmarshal(body, &data)
		i.NoError(err)

		msg, _ := data["message"].(string)
		errMsg, _ := data["error"].(string)

		i.Equal("forbidden", msg)
		i.Contains(errMsg, "user usage limit exceeded")
	})
}

func (i *integrationTestSuite) TestRateLimiterValidationFail() {
	i.Run("fail - no user-id header", func() {
		req := i.newHTTPRequest("", "", nil)

		res, err := i.server.Test(req)
		i.NoError(err)
		i.Equal(http.StatusBadRequest, res.StatusCode)

		body, err := io.ReadAll(res.Body)
		i.NoError(err)

		data := make(map[string]any)
		err = json.Unmarshal(body, &data)
		i.NoError(err)

		msg, _ := data["message"].(string)
		errMsg, _ := data["error"].(string)

		i.Equal("validation error", msg)
		i.Equal("user-id header is required", errMsg)
	})
}

func (i *integrationTestSuite) TestRateLimiterFail() {
	i.Run("fail - rate limit exceeded", func() {
		req := i.newHTTPRequest("123456", "", nil)

		res, err := i.server.Test(req, timeoutMiliseconds)
		i.NoError(err)
		i.Equal(http.StatusTooManyRequests, res.StatusCode)

		body, err := io.ReadAll(res.Body)
		i.NoError(err)

		data := make(map[string]any)
		err = json.Unmarshal(body, &data)
		i.NoError(err)

		msg, _ := data["message"].(string)
		errMsg, _ := data["error"].(string)

		i.Equal("too many requests", msg)
		i.Contains(errMsg, "too many requests")
	})
}

func (i *integrationTestSuite) TestUploadSuccess() {
	i.Run("success", func() {
		req := i.newHTTPRequest("123456", "abcdef", []byte("data"))

		res, err := i.server.Test(req, timeoutMiliseconds)
		i.NoError(err)
		i.Equal(http.StatusOK, res.StatusCode)

		body, err := io.ReadAll(res.Body)
		i.NoError(err)

		data := make(map[string]any)
		err = json.Unmarshal(body, &data)
		i.NoError(err)

		msg, _ := data["message"].(string)
		errMsg, _ := data["error"].(string)

		i.Equal("success", msg)
		i.Equal(errMsg, "")
	})
}

// Helpers
func (i *integrationTestSuite) newHTTPRequest(userID string, fileID string, fileData []byte) *http.Request {
	var req *http.Request

	if fileData != nil {
		body := bytes.Buffer{}
		mw := multipart.NewWriter(&body)

		// Add file
		fw, _ := mw.CreateFormFile("file", "testfile")
		fw.Write(fileData)

		mw.WriteField("file_id", fileID)

		mw.Close()

		req = httptest.NewRequest(http.MethodPost, "/upload", &body)

		req.Header.Set("Content-Type", mw.FormDataContentType())
	} else {
		req = httptest.NewRequest(http.MethodPost, "/upload", nil)
	}

	if userID != "" {
		req.Header.Set("user-id", userID)
	}

	return req
}
