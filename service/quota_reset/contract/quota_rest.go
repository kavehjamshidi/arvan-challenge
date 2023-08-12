package contract

import "context"

type QuotaResetService interface {
	ResetUserQuota(ctx context.Context) error
}
