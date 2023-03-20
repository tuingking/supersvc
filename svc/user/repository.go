package user

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/tuingking/supersvc/entity"
	"github.com/tuingking/supersvc/pkg/ctxkey"
	"github.com/tuingking/supersvc/pkg/mysql"
	"github.com/tuingking/supersvc/pkg/qbuilder"
)

type Repository interface {
	FindAll(ctx context.Context, p GetUserParam) ([]User, entity.Pagination, error)
	Create(ctx context.Context, v User) error
}

type repository struct {
	opt RepositoryOption
	db  mysql.MySQL
}

type RepositoryOption struct{}

func NewRepository(opt RepositoryOption, db mysql.MySQL) Repository {
	return &repository{
		opt: opt,
		db:  db,
	}
}

func (r *repository) Create(ctx context.Context, v User) error {
	logger := ctx.Value(ctxkey.ZeroLogSubLogger).(zerolog.Logger)

	res, err := r.db.Get().ExecContext(ctx, createUserQuery,
		v.ID,
		v.Name,
		v.Phone,
		v.Email,
		v.Status,
		v.CreatedAt,
	)
	if err != nil {
		logger.Err(err).Msg("failed: db.ExecContext")
		return err
	}

	rowAffected, _ := res.RowsAffected()
	logger.Debug().Int64("rows_affected", rowAffected)

	return nil
}

func (r *repository) FindAll(ctx context.Context, p GetUserParam) ([]User, entity.Pagination, error) {
	logger := ctx.Value(ctxkey.ZeroLogSubLogger).(zerolog.Logger)

	var (
		results    []User
		pagination entity.Pagination
	)

	p.Page, p.Limit = qbuilder.ValidatePageAndLimit(p.Page, p.Limit)

	qb := qbuilder.New(qbuilder.WithExtraLimit())
	clause, args, err := qb.Build(&p)
	if err != nil {
		logger.Error().Err(err).Msg("failed: qbuilder.Build")
		return results, pagination, err
	}

	rows, err := r.db.Get().QueryContext(ctx, getUserQuery+clause, args...)
	if err != nil {
		logger.Error().Err(err).Msg("failed: db.QueryContext")
		return results, pagination, err
	}

	for rows.Next() {
		var usr User
		if err := rows.Scan(
			&usr.ID,
			&usr.Name,
			&usr.Phone,
			&usr.Email,
			&usr.Status,
			&usr.CreatedAt,
		); err != nil {
			logger.Error().Err(err).Msg("failed: rows.Scan")
			return results, pagination, err
		}
		results = append(results, usr)
	}

	clausec, argsc, err := qb.BuildCount()
	if err != nil {
		logger.Error().Err(err).Msg("failed: qbuilder.BuildCount")
		return results, pagination, err
	}

	var totalData int
	row := r.db.Get().QueryRowContext(ctx, countUserQuery+clausec, argsc...)
	row.Scan(&totalData)

	size := int64(len(results))
	hasNext := len(results) > int(p.Limit)
	if hasNext {
		results = results[:p.Limit]
		size = p.Limit
	}

	pagination = entity.Pagination{
		Page:    p.Page,
		Size:    size,
		HasNext: hasNext,
		Total:   int64(totalData),
	}

	return results, pagination, nil
}
