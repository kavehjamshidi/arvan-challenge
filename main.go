package main

import (
	"github.com/kavehjamshidi/arvan-challenge/api"
	"github.com/kavehjamshidi/arvan-challenge/api/handler"
	"github.com/kavehjamshidi/arvan-challenge/driver/db/postgres"
	"github.com/kavehjamshidi/arvan-challenge/driver/db/redis"
	"github.com/kavehjamshidi/arvan-challenge/external/rate_limiter"
	"github.com/kavehjamshidi/arvan-challenge/queue"
	"github.com/kavehjamshidi/arvan-challenge/repository/file"
	"github.com/kavehjamshidi/arvan-challenge/repository/user"
	"github.com/kavehjamshidi/arvan-challenge/scheduler"
	"github.com/kavehjamshidi/arvan-challenge/service/quota_reset"
	"github.com/kavehjamshidi/arvan-challenge/service/rate_limit"
	"github.com/kavehjamshidi/arvan-challenge/service/upload"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	viper.AutomaticEnv()
	viper.ReadInConfig()

	viper.SetDefault("ENV", "dev")
	viper.SetDefault("SERVER_ADDRESS", ":4000")
	viper.SetDefault("DB_URI", "postgres://postgres:very-secret@127.0.0.1:5432/arvan?sslmode=disable")
	viper.SetDefault("REDIS_ADDRESS", "127.0.0.1:6379")
	viper.SetDefault("DB_MIGRATION_TABLE", "migrations")

	pg := postgres.Setup()
	defer pg.Close()

	postgres.Migrate(pg)

	env := viper.GetString("ENV")
	if env == "dev" {
		postgres.Seed(pg)
	}

	redisClient := redis.Setup()
	defer redisClient.Close()

	rateLimiter := rate_limiter.NewRedisRateLimiter(redisClient)
	userRepo := user.NewPostgresUserRepository(pg)
	fileRepo := file.NewRedisFileRepository(redisClient)
	noOpQueue := queue.NewNoOpQueue()

	rateLimitSVC := rate_limit.NewRateLimitService(userRepo, rateLimiter)
	uploadSVC := upload.NewUploadService(fileRepo, userRepo, noOpQueue)
	quotaResetSVC := quota_reset.NewQuotaResetService(userRepo)

	rateLimitHandler := handler.NewRateLimitHandler(rateLimitSVC)
	uploadHandler := handler.NewUploadHandler(uploadSVC)

	server := api.SetupServer(rateLimitHandler, uploadHandler)

	go func() {
		address := viper.GetString("SERVER_ADDRESS")
		if err := server.Listen(address); err != nil {
			log.Panicf("could not initialize server: %v\n", err)
		}
	}()

	go scheduler.Schedule(quotaResetSVC)

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	server.Shutdown()
}
