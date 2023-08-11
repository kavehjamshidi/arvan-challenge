package file

import (
	"context"
	"fmt"
	"github.com/kavehjamshidi/arvan-challenge/repository/file/contract"
	"github.com/kavehjamshidi/arvan-challenge/utils"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

const (
	prefix = "file_id"
)

type redisFileRepo struct {
	client *redis.Client
}

func NewRedisFileRepository(client *redis.Client) contract.FileRepository {
	return redisFileRepo{client: client}
}

func (r redisFileRepo) CheckAndStoreUniqueID(ctx context.Context, id string, ttlSeconds int) error {
	args := redis.SetArgs{
		Mode: "NX",
		TTL:  time.Duration(ttlSeconds) * time.Second,
	}

	err := r.client.SetArgs(ctx, r.generateKey(id, prefix), 1, args).Err()
	if err != nil {
		log.Println(err)
		return errors.Wrap(utils.ErrFileIDAlreadyExists, "CheckAndStoreUniqueID")
	}

	return nil
}

func (r redisFileRepo) DeleteID(ctx context.Context, id string) error {
	return r.client.Del(ctx, r.generateKey(id, prefix)).Err()
}

func (r redisFileRepo) generateKey(id, prefix string) string {
	return fmt.Sprintf("%s:%s", prefix, id)
}
