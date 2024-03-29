package transaction

import (
	"context"

	"github.com/e-fish/api/pkg/domain/transaction/model"
	"github.com/google/uuid"
)

type Repo interface {
	NewCommand(ctx context.Context) Command
	NewQuery() Query
}

type Command interface {
	CreateOrder(ctx context.Context, input model.CreateOrderInput) (*uuid.UUID, error)
	UpdateCancelOrder(ctx context.Context, input uuid.UUID) (*uuid.UUID, error)
	UpdateSuccesOrder(ctx context.Context, input uuid.UUID) (*uuid.UUID, error)

	Rollback(ctx context.Context) error
	Commit(ctx context.Context) error
}

type Query interface {
	ReadOrder(ctx context.Context, input model.ReadInput) (*model.OrderOutputPagination, error)
	ReadOrderByStatus(ctx context.Context, input model.ReadInput, status string) (*model.OrderOutputPagination, error)
	ReadOrderByID(ctx context.Context, id uuid.UUID) (*model.OrderOutput, error)
	ReadAllOrderActive(ctx context.Context) ([]*model.Order, error)

	lock() Query
}
