package repository

import (
	"context"

	dbent "github.com/meitianwang/fast-frame/ent"
	"github.com/meitianwang/fast-frame/ent/proxyproduct"
	"github.com/meitianwang/fast-frame/internal/service"
)

type proxyProductRepository struct {
	client *dbent.Client
}

func NewProxyProductRepository(client *dbent.Client) service.ProxyProductRepository {
	return &proxyProductRepository{client: client}
}

func (r *proxyProductRepository) Create(ctx context.Context, p *service.ProxyProduct) error {
	client := clientFromContext(ctx, r.client)
	created, err := client.ProxyProduct.Create().
		SetName(p.Name).
		SetDescription(p.Description).
		SetDurationDays(p.DurationDays).
		SetTrafficLimitGB(p.TrafficLimitGB).
		SetPrice(p.Price).
		SetSortOrder(p.SortOrder).
		SetIsActive(p.IsActive).
		Save(ctx)
	if err != nil {
		return translatePersistenceError(err, nil, nil)
	}
	applyProxyProductEntity(p, created)
	return nil
}

func (r *proxyProductRepository) GetByID(ctx context.Context, id int64) (*service.ProxyProduct, error) {
	client := clientFromContext(ctx, r.client)
	e, err := client.ProxyProduct.Get(ctx, id)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrProxyProductNotFound, nil)
	}
	p := &service.ProxyProduct{}
	applyProxyProductEntity(p, e)
	return p, nil
}

func (r *proxyProductRepository) Update(ctx context.Context, p *service.ProxyProduct) error {
	client := clientFromContext(ctx, r.client)
	_, err := client.ProxyProduct.UpdateOneID(p.ID).
		SetName(p.Name).
		SetDescription(p.Description).
		SetDurationDays(p.DurationDays).
		SetTrafficLimitGB(p.TrafficLimitGB).
		SetPrice(p.Price).
		SetSortOrder(p.SortOrder).
		SetIsActive(p.IsActive).
		Save(ctx)
	return translatePersistenceError(err, service.ErrProxyProductNotFound, nil)
}

func (r *proxyProductRepository) Delete(ctx context.Context, id int64) error {
	client := clientFromContext(ctx, r.client)
	err := client.ProxyProduct.DeleteOneID(id).Exec(ctx)
	return translatePersistenceError(err, service.ErrProxyProductNotFound, nil)
}

func (r *proxyProductRepository) List(ctx context.Context, activeOnly bool) ([]service.ProxyProduct, error) {
	client := clientFromContext(ctx, r.client)
	q := client.ProxyProduct.Query()
	if activeOnly {
		q = q.Where(proxyproduct.IsActive(true))
	}

	entities, err := q.Order(dbent.Asc(proxyproduct.FieldSortOrder), dbent.Asc(proxyproduct.FieldID)).All(ctx)
	if err != nil {
		return nil, err
	}

	products := make([]service.ProxyProduct, len(entities))
	for i, e := range entities {
		applyProxyProductEntity(&products[i], e)
	}
	return products, nil
}

func applyProxyProductEntity(p *service.ProxyProduct, e *dbent.ProxyProduct) {
	p.ID = e.ID
	p.Name = e.Name
	p.Description = e.Description
	p.DurationDays = e.DurationDays
	p.TrafficLimitGB = e.TrafficLimitGB
	p.Price = e.Price
	p.SortOrder = e.SortOrder
	p.IsActive = e.IsActive
	p.CreatedAt = e.CreatedAt
	p.UpdatedAt = e.UpdatedAt
}
