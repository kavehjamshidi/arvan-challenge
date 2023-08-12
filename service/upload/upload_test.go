package upload

import (
	"context"
	"github.com/kavehjamshidi/arvan-challenge/domain"
	"github.com/kavehjamshidi/arvan-challenge/mocks"
	"github.com/kavehjamshidi/arvan-challenge/utils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type uploadTestSuite struct {
	suite.Suite
	fileRepo *mocks.FileRepository
	queue    *mocks.FileUploadQueue
	userRepo *mocks.UserRepository
}

func TestUploadSuite(t *testing.T) {
	suite.Run(t, &uploadTestSuite{})
}

func (u *uploadTestSuite) SetupSubTest() {
	mockFileRepo := mocks.NewFileRepository(u.T())
	mockQueue := mocks.NewFileUploadQueue(u.T())
	mockUserRepo := mocks.NewUserRepository(u.T())

	u.fileRepo = mockFileRepo
	u.queue = mockQueue
	u.userRepo = mockUserRepo
}

func (u *uploadTestSuite) TestUploadFile() {
	file := domain.File{
		Data:   nil,
		Size:   500,
		FileID: "abc",
		UserID: "123456",
	}

	u.Run("failed - could not check fileID uniqueness", func() {
		uploadSVC := NewUploadService(u.fileRepo, u.userRepo, u.queue)

		u.fileRepo.On("CheckAndStoreUniqueID", context.TODO(), file.FileID, mock.Anything).
			Return(errors.New("unknown error"))

		err := uploadSVC.UploadFile(context.TODO(), file)

		u.ErrorContains(err, "unknown error")
	})

	u.Run("failed - duplicate fileID", func() {
		uploadSVC := NewUploadService(u.fileRepo, u.userRepo, u.queue)

		u.fileRepo.On("CheckAndStoreUniqueID", context.TODO(), file.FileID, mock.Anything).
			Return(utils.ErrFileIDAlreadyExists)

		err := uploadSVC.UploadFile(context.TODO(), file)

		u.ErrorIs(err, utils.ErrFileIDAlreadyExists)
	})

	u.Run("failed - could not increment usage", func() {
		uploadSVC := NewUploadService(u.fileRepo, u.userRepo, u.queue)

		u.fileRepo.On("CheckAndStoreUniqueID", context.TODO(), file.FileID, mock.Anything).
			Return(nil)

		u.userRepo.On("IncrementUsage", context.TODO(), file.UserID, file.Size).
			Return(errors.New("unknown error"))

		u.fileRepo.On("DeleteID", context.TODO(), file.FileID).Return(nil)

		err := uploadSVC.UploadFile(context.TODO(), file)

		u.ErrorContains(err, "unknown error")
	})

	u.Run("failed - usage quota exceeded", func() {
		uploadSVC := NewUploadService(u.fileRepo, u.userRepo, u.queue)

		u.fileRepo.On("CheckAndStoreUniqueID", context.TODO(), file.FileID, mock.Anything).
			Return(nil)

		u.userRepo.On("IncrementUsage", context.TODO(), file.UserID, file.Size).
			Return(utils.ErrUsageLimitExceeded)

		u.fileRepo.On("DeleteID", context.TODO(), file.FileID).Return(nil)

		err := uploadSVC.UploadFile(context.TODO(), file)

		u.ErrorIs(err, utils.ErrUsageLimitExceeded)
	})

	u.Run("failed - queue error", func() {
		uploadSVC := NewUploadService(u.fileRepo, u.userRepo, u.queue)

		u.fileRepo.On("CheckAndStoreUniqueID", context.TODO(), file.FileID, mock.Anything).
			Return(nil)

		u.userRepo.On("IncrementUsage", context.TODO(), file.UserID, file.Size).
			Return(nil)

		u.queue.On("Store", file).Return(errors.New("unknown error"))

		u.userRepo.On("DecrementUsage", context.TODO(), file.UserID, file.Size).
			Return(nil)

		u.fileRepo.On("DeleteID", context.TODO(), file.FileID).Return(nil)

		err := uploadSVC.UploadFile(context.TODO(), file)

		u.ErrorContains(err, "unknown error")
	})

	u.Run("success", func() {
		uploadSVC := NewUploadService(u.fileRepo, u.userRepo, u.queue)

		u.fileRepo.On("CheckAndStoreUniqueID", context.TODO(), file.FileID, mock.Anything).
			Return(nil)

		u.userRepo.On("IncrementUsage", context.TODO(), file.UserID, file.Size).
			Return(nil)

		u.queue.On("Store", file).Return(nil)

		err := uploadSVC.UploadFile(context.TODO(), file)

		u.NoError(err)
	})
}
