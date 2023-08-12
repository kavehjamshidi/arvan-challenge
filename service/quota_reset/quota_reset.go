package quota_reset

import (
	"context"
	"github.com/kavehjamshidi/arvan-challenge/repository/user/contract"
	contract2 "github.com/kavehjamshidi/arvan-challenge/service/quota_reset/contract"
	"github.com/pkg/errors"
	"time"
)

type quotaResetService struct {
	userRepo contract.UserRepository
}

func NewQuotaResetService(userRepo contract.UserRepository) contract2.QuotaResetService {
	return quotaResetService{userRepo: userRepo}
}

func (q quotaResetService) ResetUserQuota(ctx context.Context) error {
	end := time.Now()

	err := q.userRepo.ResetUsage(ctx, end)
	if err != nil {
		return errors.Wrap(err, "ResetUserQuota")
	}

	return nil
}
