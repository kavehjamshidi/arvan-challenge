package upload

import (
	"context"
	"github.com/kavehjamshidi/arvan-challenge/domain"
	contract3 "github.com/kavehjamshidi/arvan-challenge/queue/contract"
	"github.com/kavehjamshidi/arvan-challenge/repository/file/contract"
	contract2 "github.com/kavehjamshidi/arvan-challenge/repository/user/contract"
	contract4 "github.com/kavehjamshidi/arvan-challenge/service/upload/contract"
	"github.com/kavehjamshidi/arvan-challenge/utils"
	"github.com/pkg/errors"
)

const (
	fileIDTTLSeconds = 30 * 24 * 3600
)

type uploadService struct {
	fileRepo contract.FileRepository
	userRepo contract2.UserRepository
	queue    contract3.FileUploadQueue
}

func NewUploadService(
	fileRepo contract.FileRepository,
	userRepo contract2.UserRepository,
	queue contract3.FileUploadQueue,
) contract4.UploadService {
	return uploadService{
		fileRepo: fileRepo,
		userRepo: userRepo,
		queue:    queue,
	}
}

func (us uploadService) UploadFile(ctx context.Context, file domain.File) error {
	err := us.fileRepo.CheckAndStoreUniqueID(ctx, file.FileID, fileIDTTLSeconds)
	if err != nil {
		if errors.Is(err, utils.ErrFileIDAlreadyExists) {
			return errors.Wrap(err, "UploadFile")
		}
		return errors.Wrap(err, "UploadFile")
	}

	err = us.userRepo.IncrementUsage(ctx, file.UserID, file.Size)
	if err != nil {
		us.fileRepo.DeleteID(ctx, file.FileID)
		if errors.Is(err, utils.ErrUsageLimitExceeded) {
			return errors.Wrap(err, "UploadFile")
		}
		return errors.Wrap(err, "UploadFile")
	}

	err = us.queue.Store(file)
	if err != nil {
		us.fileRepo.DeleteID(ctx, file.FileID)
		us.userRepo.DecrementUsage(ctx, file.UserID, file.Size)
		return errors.Wrap(err, "UploadFile")
	}

	return nil
}
