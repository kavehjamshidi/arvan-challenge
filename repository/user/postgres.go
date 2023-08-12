package user

import (
	"context"
	"database/sql"
	"github.com/kavehjamshidi/arvan-challenge/repository/user/contract"
	"github.com/kavehjamshidi/arvan-challenge/utils"
	"github.com/pkg/errors"
	"log"
	"time"
)

type pgUserRepo struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) contract.UserRepository {
	return pgUserRepo{db: db}
}

func (p pgUserRepo) ResetUsage(ctx context.Context, end time.Time) error {
	query := `UPDATE user_usage SET quota_usage = 0, start_date = end_date,
                      end_date = end_date + INTERVAL '1 month'
WHERE end_date <= $1;`

	_, err := p.db.ExecContext(ctx, query, end)
	if err != nil {
		log.Println(err)
		return errors.Wrap(err, "ResetUserQuota")
	}

	return nil
}

func (p pgUserRepo) GetUserRateLimit(ctx context.Context, userID string) (int, error) {
	query := `SELECT rate_limit FROM users WHERE id = $1`

	row := p.db.QueryRowContext(ctx, query, userID)
	if row.Err() != nil {
		return 0, errors.Wrap(row.Err(), "GetUserRateLimit")
	}

	var limit int
	err := row.Scan(&limit)
	if err != nil {
		return 0, errors.Wrap(err, "GetUserRateLimit")
	}

	return limit, nil
}

func (p pgUserRepo) IncrementUsage(ctx context.Context, userID string, usage int64) error {
	query := `UPDATE user_usage SET quota_usage = quota_usage + $1 
WHERE user_id = $2 AND quota - quota_usage >= $3`

	res, err := p.db.ExecContext(ctx, query, usage, userID, usage)
	if err != nil {
		log.Println(err)
		return errors.Wrap(err, "IncrementUsage")
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Println(err)
		return errors.Wrap(err, "IncrementUsage")
	}

	if rowsAffected == 0 {
		return errors.Wrap(utils.ErrUsageLimitExceeded, "IncrementUsage")
	}

	return nil
}

func (p pgUserRepo) DecrementUsage(ctx context.Context, userID string, usage int64) error {
	query := `UPDATE user_usage SET quota_usage = quota_usage - $1 
WHERE user_id = $2`

	_, err := p.db.ExecContext(ctx, query, usage, userID, usage)
	if err != nil {
		log.Println(err)
		return errors.Wrap(err, "DecrementUsage")
	}

	return nil
}
