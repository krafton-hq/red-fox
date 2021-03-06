package service_decorator

import (
	"context"

	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	"github.com/krafton-hq/red-fox/server/pkg/transactional"
	"github.com/krafton-hq/red-fox/server/services/services"
)

type TransactionalClusterService[T any] struct {
	service services.ClusterService[T]

	tr transactional.Transactional
}

func NewTransactionalClusterService[T any](service services.ClusterService[T], tr transactional.Transactional) services.ClusterService[T] {
	return &TransactionalClusterService[T]{service: service, tr: tr}
}

func (s *TransactionalClusterService[T]) Init(ctx context.Context) error {
	ctx = s.tr.WithDatabaseContext(ctx)

	err := s.tr.WithTransaction(ctx, func(ctx context.Context, tx *sqlx.Tx, dialect goqu.SQLDialect) error {
		return s.service.Init(ctx)
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *TransactionalClusterService[T]) Get(ctx context.Context, name string) (T, error) {
	ctx = s.tr.WithDatabaseContext(ctx)

	var result T
	var err error
	err = s.tr.WithTransaction(ctx, func(ctx context.Context, tx *sqlx.Tx, dialect goqu.SQLDialect) error {
		result, err = s.service.Get(ctx, name)
		return err
	})
	if err != nil {
		return result, err
	}
	return result, nil
}

func (s *TransactionalClusterService[T]) List(ctx context.Context, labelSelectors map[string]string) ([]T, error) {
	ctx = s.tr.WithDatabaseContext(ctx)

	var result []T
	var err error
	err = s.tr.WithTransaction(ctx, func(ctx context.Context, tx *sqlx.Tx, dialect goqu.SQLDialect) error {
		result, err = s.service.List(ctx, labelSelectors)
		return err
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *TransactionalClusterService[T]) Create(ctx context.Context, obj T) error {
	ctx = s.tr.WithDatabaseContext(ctx)

	err := s.tr.WithTransaction(ctx, func(ctx context.Context, tx *sqlx.Tx, dialect goqu.SQLDialect) error {
		return s.service.Create(ctx, obj)
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *TransactionalClusterService[T]) Update(ctx context.Context, obj T) error {
	ctx = s.tr.WithDatabaseContext(ctx)

	err := s.tr.WithTransaction(ctx, func(ctx context.Context, tx *sqlx.Tx, dialect goqu.SQLDialect) error {
		return s.service.Update(ctx, obj)
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *TransactionalClusterService[T]) Delete(ctx context.Context, name string) error {
	ctx = s.tr.WithDatabaseContext(ctx)

	err := s.tr.WithTransaction(ctx, func(ctx context.Context, tx *sqlx.Tx, dialect goqu.SQLDialect) error {
		return s.service.Delete(ctx, name)
	})
	if err != nil {
		return err
	}
	return nil
}
