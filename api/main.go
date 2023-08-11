package main

import (
	"github.com/kavehjamshidi/arvan-challenge/api/http"
	handler2 "github.com/kavehjamshidi/arvan-challenge/api/http/handler"
	"github.com/kavehjamshidi/arvan-challenge/driver/db/postgres"
	"github.com/kavehjamshidi/arvan-challenge/driver/db/redis"
	"github.com/kavehjamshidi/arvan-challenge/external/rate_limiter"
	"github.com/kavehjamshidi/arvan-challenge/queue"
	"github.com/kavehjamshidi/arvan-challenge/repository/file"
	"github.com/kavehjamshidi/arvan-challenge/repository/user"
	"github.com/kavehjamshidi/arvan-challenge/service/rate_limit"
	"github.com/kavehjamshidi/arvan-challenge/service/upload"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	viper.AutomaticEnv()
	viper.ReadInConfig()

	viper.SetDefault("SERVER_ADDRESS", ":4000")
	viper.SetDefault("DB_URI", "postgres://postgres:very-secret@127.0.0.1:5432/arvan?sslmode=disable")
	viper.SetDefault("REDIS_ADDRESS", "127.0.0.1:6379")
	viper.SetDefault("DB_MIGRATION_TABLE", "migrations")

	pg := postgres.Setup()
	defer pg.Close()

	postgres.Migrate(pg)

	redisClient := redis.Setup()
	defer redisClient.Close()

	rateLimiter := rate_limiter.NewRedisRateLimiter(redisClient)
	userRepo := user.NewPostgresUserRepository(pg)
	fileRepo := file.NewRedisFileRepository(redisClient)
	noOpQueue := queue.NewNoOpQueue()

	rateLimitSVC := rate_limit.NewRateLimitService(userRepo, rateLimiter)
	uploadSVC := upload.NewUploadService(fileRepo, userRepo, noOpQueue)

	rateLimitHandler := handler2.NewRateLimitHandler(rateLimitSVC)
	uploadHandler := handler2.NewUploadHandler(uploadSVC)

	go http.SetupServer(rateLimitHandler, uploadHandler)

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
