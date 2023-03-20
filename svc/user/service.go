package user

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/tuingking/supersvc/entity"
	"github.com/tuingking/supersvc/pkg/logger"
)

type Service interface {
	GetUser(ctx context.Context, p GetUserParam) ([]User, entity.Pagination, error)
	CreateUser(ctx context.Context, v User) (User, error)
}

type service struct {
	opt  Option
	repo Repository
}

type Option struct {
	Repository RepositoryOption
}

func NewService(opt Option, repo Repository) Service {
	return &service{
		opt:  opt,
		repo: repo,
	}
}

func (s *service) GetUser(ctx context.Context, p GetUserParam) ([]User, entity.Pagination, error) {
	return s.repo.FindAll(ctx, p)
}

func (s *service) CreateUser(ctx context.Context, v User) (User, error) {
	log := logger.Get(ctx)

	// create user id
	v.ID = uuid.New().String()
	v.CreatedAt = time.Now()

	if err := s.repo.Create(ctx, v); err != nil {
		log.Err(err).Msg("failed: create user")
		return v, err
	}

	log.Debug().Str("user_id", v.ID).Msg("user created")

	return v, nil
}
