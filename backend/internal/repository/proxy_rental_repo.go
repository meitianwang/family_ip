package repository

import (
	"context"
	"time"

	dbent "github.com/meitianwang/fast-frame/ent"
	"github.com/meitianwang/fast-frame/ent/proxyrental"
	"github.com/meitianwang/fast-frame/internal/domain"
	"github.com/meitianwang/fast-frame/internal/pkg/pagination"
	"github.com/meitianwang/fast-frame/internal/service"
)

type proxyRentalRepository struct {
	client *dbent.Client
}

func NewProxyRentalRepository(client *dbent.Client) service.ProxyRentalRepository {
	return &proxyRentalRepository{client: client}
}

func (r *proxyRentalRepository) Create(ctx context.Context, rental *service.ProxyRental) error {
	client := clientFromContext(ctx, r.client)
	q := client.ProxyRental.Create().
		SetUserID(rental.UserID).
		SetNodeID(rental.NodeID).
		SetProductID(rental.ProductID).
		SetStatus(rental.Status).
		SetTrafficUsedBytes(rental.TrafficUsedBytes).
		SetTrafficLimitBytes(rental.TrafficLimitBytes).
		SetNillablePaymentOrderID(rental.PaymentOrderID)

	created, err := q.Save(ctx)
	if err != nil {
		return translatePersistenceError(err, nil, nil)
	}
	applyProxyRentalEntity(rental, created)
	return nil
}

func (r *proxyRentalRepository) GetByID(ctx context.Context, id int64) (*service.ProxyRental, error) {
	client := clientFromContext(ctx, r.client)
	e, err := client.ProxyRental.Query().
		Where(proxyrental.ID(id)).
		WithNode().
		WithProduct().
		Only(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrProxyRentalNotFound, nil)
	}
	rental := &service.ProxyRental{}
	applyProxyRentalEntity(rental, e)
	return rental, nil
}

func (r *proxyRentalRepository) GetByPaymentOrderID(ctx context.Context, orderID int64) (*service.ProxyRental, error) {
	client := clientFromContext(ctx, r.client)
	e, err := client.ProxyRental.Query().
		Where(proxyrental.PaymentOrderID(orderID)).
		Only(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrProxyRentalNotFound, nil)
	}
	rental := &service.ProxyRental{}
	applyProxyRentalEntity(rental, e)
	return rental, nil
}

func (r *proxyRentalRepository) Update(ctx context.Context, rental *service.ProxyRental) error {
	client := clientFromContext(ctx, r.client)
	u := client.ProxyRental.UpdateOneID(rental.ID).
		SetStatus(rental.Status).
		SetTrafficUsedBytes(rental.TrafficUsedBytes).
		SetNillablePaymentOrderID(rental.PaymentOrderID)

	if rental.StartedAt != nil {
		u = u.SetStartedAt(*rental.StartedAt)
	}
	if rental.ExpiresAt != nil {
		u = u.SetExpiresAt(*rental.ExpiresAt)
	}

	_, err := u.Save(ctx)
	return translatePersistenceError(err, service.ErrProxyRentalNotFound, nil)
}

func (r *proxyRentalRepository) List(ctx context.Context, filter service.ProxyRentalFilter, params pagination.PaginationParams) ([]service.ProxyRental, *pagination.PaginationResult, error) {
	client := clientFromContext(ctx, r.client)
	q := client.ProxyRental.Query().WithNode().WithProduct()

	if filter.UserID != nil {
		q = q.Where(proxyrental.UserID(*filter.UserID))
	}
	if filter.NodeID != nil {
		q = q.Where(proxyrental.NodeID(*filter.NodeID))
	}
	if filter.Status != "" {
		q = q.Where(proxyrental.Status(filter.Status))
	}

	total, err := q.Count(ctx)
	if err != nil {
		return nil, nil, err
	}

	offset := (params.Page - 1) * params.PageSize
	entities, err := q.Order(dbent.Desc(proxyrental.FieldCreatedAt)).
		Offset(offset).
		Limit(params.PageSize).
		All(ctx)
	if err != nil {
		return nil, nil, err
	}

	rentals := make([]service.ProxyRental, len(entities))
	for i, e := range entities {
		applyProxyRentalEntity(&rentals[i], e)
	}
	pages := total / params.PageSize
	if total%params.PageSize > 0 {
		pages++
	}
	result := &pagination.PaginationResult{Total: int64(total), Page: params.Page, PageSize: params.PageSize, Pages: pages}
	return rentals, result, nil
}

func (r *proxyRentalRepository) ListByUserID(ctx context.Context, userID int64, params pagination.PaginationParams) ([]service.ProxyRental, *pagination.PaginationResult, error) {
	return r.List(ctx, service.ProxyRentalFilter{UserID: &userID}, params)
}

func (r *proxyRentalRepository) ListExpiredActive(ctx context.Context, now time.Time, limit int) ([]service.ProxyRental, error) {
	client := clientFromContext(ctx, r.client)
	entities, err := client.ProxyRental.Query().
		Where(
			proxyrental.StatusEQ(domain.ProxyRentalStatusActive),
			proxyrental.ExpiresAtLT(now),
		).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, err
	}

	rentals := make([]service.ProxyRental, len(entities))
	for i, e := range entities {
		applyProxyRentalEntity(&rentals[i], e)
	}
	return rentals, nil
}

func applyProxyRentalEntity(rental *service.ProxyRental, e *dbent.ProxyRental) {
	rental.ID = e.ID
	rental.UserID = e.UserID
	rental.NodeID = e.NodeID
	rental.ProductID = e.ProductID
	rental.PaymentOrderID = e.PaymentOrderID
	rental.Status = e.Status
	rental.StartedAt = e.StartedAt
	rental.ExpiresAt = e.ExpiresAt
	rental.TrafficUsedBytes = e.TrafficUsedBytes
	rental.TrafficLimitBytes = e.TrafficLimitBytes
	rental.CreatedAt = e.CreatedAt
	rental.UpdatedAt = e.UpdatedAt

	if e.Edges.Node != nil {
		node := &service.ProxyNode{}
		applyProxyNodeEntity(node, e.Edges.Node)
		rental.Node = node
	}
	if e.Edges.Product != nil {
		product := &service.ProxyProduct{}
		applyProxyProductEntity(product, e.Edges.Product)
		rental.Product = product
	}
}
