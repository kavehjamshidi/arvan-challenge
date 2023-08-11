package contract

import "context"

type FileRepository interface {
	CheckAndStoreUniqueID(ctx context.Context, id string, ttlSeconds int) (err error)
	DeleteID(ctx context.Context, id string) (err error)
}
