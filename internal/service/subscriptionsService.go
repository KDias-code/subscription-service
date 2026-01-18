package service

import (
	"context"
	"github.com/google/uuid"
	"subscription-service/internal/model"
	"subscription-service/internal/repository"
)

type ISubscriptionService interface {
	FindByID(ctx context.Context, id uuid.UUID) (*model.Subscriptions, error)
	Create(ctx context.Context, data model.Subscriptions) error
	Update(ctx context.Context, data model.Subscriptions) error
	Delete(ctx context.Context, id uuid.UUID) error
	SubscriptionsSum(ctx context.Context, data model.SubscriptionsSum) (uint64, error)
}

type SubscriptionsService struct {
	repo repository.ISubscriptionsRepository
}

func NewSubscriptionsService(repo repository.ISubscriptionsRepository) ISubscriptionService {
	return &SubscriptionsService{
		repo: repo,
	}
}

func (s SubscriptionsService) FindByID(ctx context.Context, id uuid.UUID) (*model.Subscriptions, error) {
	result, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return result, err
}

func (s SubscriptionsService) Create(ctx context.Context, data model.Subscriptions) error {
	err := s.repo.Create(ctx, &data)
	if err != nil {
		return err
	}

	return nil
}

func (s SubscriptionsService) Update(ctx context.Context, data model.Subscriptions) error {
	err := s.repo.Update(ctx, &data)
	if err != nil {
		return err
	}

	return nil
}

func (s SubscriptionsService) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s SubscriptionsService) SubscriptionsSum(ctx context.Context, data model.SubscriptionsSum) (uint64, error) {
	amount, err := s.repo.SumPrice(ctx, data)
	if err != nil {
		return 0, err
	}

	return amount, nil
}
