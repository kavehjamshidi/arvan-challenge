package quota_reset

import (
	"context"
	"github.com/kavehjamshidi/arvan-challenge/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type quotaResetServiceTestSuite struct {
	suite.Suite
	userRepo *mocks.UserRepository
}

func TestQuotaResetSuite(t *testing.T) {
	suite.Run(t, &quotaResetServiceTestSuite{})
}

func (q *quotaResetServiceTestSuite) SetupSubTest() {
	mockUserRepo := mocks.NewUserRepository(q.T())

	q.userRepo = mockUserRepo
}

func (q *quotaResetServiceTestSuite) TestResetUserQuota() {
	q.Run("failed - repo error", func() {
		quotaSVC := NewQuotaResetService(q.userRepo)

		q.userRepo.On("ResetUsage", context.TODO(), mock.Anything).
			Return(errors.New("unknown error"))

		err := quotaSVC.ResetUserQuota(context.TODO())

		q.ErrorContains(err, "unknown error")
	})

	q.Run("success", func() {
		quotaSVC := NewQuotaResetService(q.userRepo)

		q.userRepo.On("ResetUsage", context.TODO(), mock.Anything).
			Return(nil)

		err := quotaSVC.ResetUserQuota(context.TODO())

		q.NoError(err)
	})

}
