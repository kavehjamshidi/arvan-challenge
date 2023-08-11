package contract

import (
	"context"
	"github.com/kavehjamshidi/arvan-challenge/domain"
)

type UploadService interface {
	UploadFile(ctx context.Context, file domain.File) (err error)
}
