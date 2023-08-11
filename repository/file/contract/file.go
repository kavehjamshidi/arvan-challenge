package contract

import (
	"context"
	"time"
)

type FileRepository interface {
	CheckAndStoreUniqueID(ctx context.Context, id string, ttl time.Duration) (err error)
	DeleteID(ctx context.Context, id string) (err error)
}
