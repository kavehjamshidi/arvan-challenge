package user

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kavehjamshidi/arvan-challenge/utils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type postgresUserRepoSuite struct {
	suite.Suite
	db   *sql.DB
	mock sqlmock.Sqlmock
}

func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, &postgresUserRepoSuite{})
}

func (p *postgresUserRepoSuite) SetupSuite() {
	db, mock, err := sqlmock.New()
	if err != nil {
		p.FailNow(err.Error())
	}

	p.db = db
	p.mock = mock
}

func (p *postgresUserRepoSuite) TearDownSuite() {
	p.db.Close()
}

func (p *postgresUserRepoSuite) TestGetUserRateLimit() {
	userID := "123456"

	p.Run("success", func() {
		repo := NewPostgresUserRepository(p.db)

		p.mock.ExpectQuery(`SELECT rate_limit FROM users WHERE id = \$1`).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"rate_limit"}).FromCSVString("100"))

		limit, err := repo.GetUserRateLimit(context.TODO(), userID)

		p.NoError(err)
		p.Equal(100, limit)

		err = p.mock.ExpectationsWereMet()
		p.NoError(err)
	})

	p.Run("failed - row error", func() {
		repo := NewPostgresUserRepository(p.db)

		p.mock.ExpectQuery(`SELECT rate_limit FROM users WHERE id = \$1`).
			WithArgs(userID).
			WillReturnError(errors.New("unknown error"))

		limit, err := repo.GetUserRateLimit(context.TODO(), userID)

		p.ErrorContains(err, "unknown error")
		p.Equal(0, limit)

		err = p.mock.ExpectationsWereMet()
		p.NoError(err)
	})

	p.Run("failed - scan error", func() {
		repo := NewPostgresUserRepository(p.db)

		p.mock.ExpectQuery(`SELECT rate_limit FROM users WHERE id = \$1`).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"rate_limit"}).FromCSVString("abc"))

		limit, err := repo.GetUserRateLimit(context.TODO(), userID)

		p.ErrorContains(err, "sql: Scan error")
		p.Equal(0, limit)

		err = p.mock.ExpectationsWereMet()
		p.NoError(err)
	})
}

func (p *postgresUserRepoSuite) TestIncrementUsage() {
	userID := "123456"
	var usage int64 = 200

	p.Run("success", func() {
		repo := NewPostgresUserRepository(p.db)

		p.mock.ExpectExec(`UPDATE user_usage SET quota_usage = quota_usage + `).
			WithArgs(usage, userID, usage).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.IncrementUsage(context.TODO(), userID, usage)

		p.NoError(err)
	})

	p.Run("fail - usage exceeded", func() {
		repo := NewPostgresUserRepository(p.db)

		p.mock.ExpectExec(`UPDATE user_usage SET quota_usage = quota_usage + `).
			WithArgs(usage, userID, usage).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.IncrementUsage(context.TODO(), userID, usage)

		p.ErrorIs(err, utils.ErrUsageLimitExceeded)
	})

	p.Run("fail - exec error", func() {
		repo := NewPostgresUserRepository(p.db)

		p.mock.ExpectExec(`UPDATE user_usage SET quota_usage = quota_usage + `).
			WithArgs(usage, userID, usage).
			WillReturnError(errors.New("unknown error"))

		err := repo.IncrementUsage(context.TODO(), userID, usage)

		p.ErrorContains(err, "unknown error")
	})

	p.Run("fail - rows affected error", func() {
		repo := NewPostgresUserRepository(p.db)

		p.mock.ExpectExec(`UPDATE user_usage SET quota_usage = quota_usage + `).
			WithArgs(usage, userID, usage).
			WillReturnResult(sqlmock.NewErrorResult(errors.New("rows affected error")))

		err := repo.IncrementUsage(context.TODO(), userID, usage)

		p.ErrorContains(err, "rows affected error")
	})
}

func (p *postgresUserRepoSuite) TestDecrementUsage() {
	userID := "123456"
	var usage int64 = 200

	p.Run("success", func() {
		repo := NewPostgresUserRepository(p.db)

		p.mock.ExpectExec(`UPDATE user_usage SET quota_usage = quota_usage - `).
			WithArgs(usage, userID, usage).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.DecrementUsage(context.TODO(), userID, usage)

		p.NoError(err)
	})

	p.Run("fail - exec error", func() {
		repo := NewPostgresUserRepository(p.db)

		p.mock.ExpectExec(`UPDATE user_usage SET quota_usage = quota_usage - `).
			WithArgs(usage, userID, usage).
			WillReturnError(errors.New("unknown error"))

		err := repo.DecrementUsage(context.TODO(), userID, usage)

		p.ErrorContains(err, "unknown error")
	})
}

func (p *postgresUserRepoSuite) TestResetUsage() {
	end := time.Now()

	p.Run("success", func() {
		repo := NewPostgresUserRepository(p.db)

		p.mock.ExpectExec(`UPDATE user_usage SET quota_usage = 0, start_date = end_date,`).
			WithArgs(end).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.ResetUsage(context.TODO(), end)

		p.NoError(err)
	})

	p.Run("fail - exec error", func() {
		repo := NewPostgresUserRepository(p.db)

		p.mock.ExpectExec(`UPDATE user_usage SET quota_usage = 0, start_date = end_date,`).
			WithArgs(end).
			WillReturnError(errors.New("unknown error"))

		err := repo.ResetUsage(context.TODO(), end)

		p.ErrorContains(err, "unknown error")
	})
}
