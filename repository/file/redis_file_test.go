package file

import (
	"context"
	"github.com/go-redis/redismock/v9"
	"github.com/kavehjamshidi/arvan-challenge/utils"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type redisFileSuite struct {
	suite.Suite
	redisClient *redis.Client
	mock        redismock.ClientMock
}

func TestRateLimiterSuite(t *testing.T) {
	suite.Run(t, &redisFileSuite{})
}

func (r *redisFileSuite) SetupSuite() {
	db, mock := redismock.NewClientMock()
	r.redisClient = db
	r.mock = mock
}

func (r *redisFileSuite) TestCheckAndStoreUniqueID() {
	ttl := 10 * time.Second
	fileID := "123456"

	r.Run("success", func() {
		repo := NewRedisFileRepository(r.redisClient)

		r.mock.ExpectSetArgs("file_id:123456", 1, redis.SetArgs{
			Mode: "NX",
			TTL:  ttl,
		}).SetVal("OK")

		err := repo.CheckAndStoreUniqueID(context.TODO(), fileID, ttl)

		r.NoError(err)
	})

	r.Run("failed", func() {
		repo := NewRedisFileRepository(r.redisClient)

		r.mock.ExpectSetArgs("file_id:123456", 1, redis.SetArgs{
			Mode: "NX",
			TTL:  ttl,
		}).SetErr(errors.New("unknown error"))

		err := repo.CheckAndStoreUniqueID(context.TODO(), fileID, ttl)

		r.ErrorIs(err, utils.ErrFileIDAlreadyExists)
	})
}

func (r *redisFileSuite) TestDeleteID() {
	fileID := "123456"

	r.Run("success", func() {
		repo := NewRedisFileRepository(r.redisClient)

		r.mock.ExpectDel("file_id:123456").SetVal(1)

		err := repo.DeleteID(context.TODO(), fileID)

		r.NoError(err)
	})

	r.Run("failed", func() {
		repo := NewRedisFileRepository(r.redisClient)

		r.mock.ExpectDel("file_id:123456").SetErr(errors.New("unknown error"))

		err := repo.DeleteID(context.TODO(), fileID)

		r.ErrorContains(err, "unknown error")
	})
}
