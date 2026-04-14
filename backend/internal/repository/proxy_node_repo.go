package repository

import (
	"context"

	dbent "github.com/meitianwang/fast-frame/ent"
	"github.com/meitianwang/fast-frame/ent/proxynode"
	"github.com/meitianwang/fast-frame/internal/pkg/pagination"
	"github.com/meitianwang/fast-frame/internal/service"
)

type proxyNodeRepository struct {
	client *dbent.Client
}

func NewProxyNodeRepository(client *dbent.Client) service.ProxyNodeRepository {
	return &proxyNodeRepository{client: client}
}

func (r *proxyNodeRepository) Create(ctx context.Context, node *service.ProxyNode) error {
	client := clientFromContext(ctx, r.client)
	created, err := client.ProxyNode.Create().
		SetIPAddress(node.IPAddress).
		SetCountry(node.Country).
		SetCountryCode(node.CountryCode).
		SetCity(node.City).
		SetIsp(node.ISP).
		SetHTTPPort(node.HTTPPort).
		SetVlessPort(node.VlessPort).
		SetVlessNetwork(node.VlessNetwork).
		SetVlessTLS(node.VlessTLS).
		SetVlessSni(node.VlessSNI).
		SetVlessWsPath(node.VlessWSPath).
		SetTags(node.Tags).
		SetStatus(node.Status).
		SetDescription(node.Description).
		Save(ctx)
	if err != nil {
		return translatePersistenceError(err, nil, nil)
	}
	applyProxyNodeEntity(node, created)
	return nil
}

func (r *proxyNodeRepository) GetByID(ctx context.Context, id int64) (*service.ProxyNode, error) {
	client := clientFromContext(ctx, r.client)
	e, err := client.ProxyNode.Get(ctx, id)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrProxyNodeNotFound, nil)
	}
	node := &service.ProxyNode{}
	applyProxyNodeEntity(node, e)
	return node, nil
}

func (r *proxyNodeRepository) Update(ctx context.Context, node *service.ProxyNode) error {
	client := clientFromContext(ctx, r.client)
	_, err := client.ProxyNode.UpdateOneID(node.ID).
		SetIPAddress(node.IPAddress).
		SetCountry(node.Country).
		SetCountryCode(node.CountryCode).
		SetCity(node.City).
		SetIsp(node.ISP).
		SetHTTPPort(node.HTTPPort).
		SetVlessPort(node.VlessPort).
		SetVlessNetwork(node.VlessNetwork).
		SetVlessTLS(node.VlessTLS).
		SetVlessSni(node.VlessSNI).
		SetVlessWsPath(node.VlessWSPath).
		SetTags(node.Tags).
		SetStatus(node.Status).
		SetDescription(node.Description).
		Save(ctx)
	return translatePersistenceError(err, service.ErrProxyNodeNotFound, nil)
}

func (r *proxyNodeRepository) SoftDelete(ctx context.Context, id int64) error {
	client := clientFromContext(ctx, r.client)
	err := client.ProxyNode.DeleteOneID(id).Exec(ctx)
	return translatePersistenceError(err, service.ErrProxyNodeNotFound, nil)
}

func (r *proxyNodeRepository) List(ctx context.Context, filter service.ProxyNodeFilter, params pagination.PaginationParams) ([]service.ProxyNode, *pagination.PaginationResult, error) {
	client := clientFromContext(ctx, r.client)
	q := client.ProxyNode.Query()
	if filter.CountryCode != "" {
		q = q.Where(proxynode.CountryCode(filter.CountryCode))
	}
	if filter.Status != "" {
		q = q.Where(proxynode.Status(filter.Status))
	}

	total, err := q.Count(ctx)
	if err != nil {
		return nil, nil, err
	}

	offset := (params.Page - 1) * params.PageSize
	entities, err := q.Order(dbent.Asc(proxynode.FieldID)).
		Offset(offset).
		Limit(params.PageSize).
		All(ctx)
	if err != nil {
		return nil, nil, err
	}

	nodes := make([]service.ProxyNode, len(entities))
	for i, e := range entities {
		applyProxyNodeEntity(&nodes[i], e)
	}
	pages := total / params.PageSize
	if total%params.PageSize > 0 {
		pages++
	}
	result := &pagination.PaginationResult{Total: int64(total), Page: params.Page, PageSize: params.PageSize, Pages: pages}
	return nodes, result, nil
}

func (r *proxyNodeRepository) ListAvailable(ctx context.Context, filter service.ProxyNodeFilter) ([]service.ProxyNode, error) {
	client := clientFromContext(ctx, r.client)
	q := client.ProxyNode.Query().Where(proxynode.StatusEQ("available"))
	if filter.CountryCode != "" {
		q = q.Where(proxynode.CountryCode(filter.CountryCode))
	}

	entities, err := q.Order(dbent.Asc(proxynode.FieldID)).All(ctx)
	if err != nil {
		return nil, err
	}

	nodes := make([]service.ProxyNode, len(entities))
	for i, e := range entities {
		applyProxyNodeEntity(&nodes[i], e)
	}
	return nodes, nil
}

func (r *proxyNodeRepository) LockForRental(ctx context.Context, id int64) (*service.ProxyNode, error) {
	// Use Ent's ForUpdate modifier to lock the row in a transaction.
	client := clientFromContext(ctx, r.client)
	e, err := client.ProxyNode.Query().
		Where(proxynode.ID(id)).
		ForUpdate().
		Only(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrProxyNodeNotFound, nil)
	}
	node := &service.ProxyNode{}
	applyProxyNodeEntity(node, e)
	return node, nil
}

func (r *proxyNodeRepository) SetStatus(ctx context.Context, id int64, status string) error {
	client := clientFromContext(ctx, r.client)
	err := client.ProxyNode.UpdateOneID(id).SetStatus(status).Exec(ctx)
	return translatePersistenceError(err, service.ErrProxyNodeNotFound, nil)
}

func applyProxyNodeEntity(node *service.ProxyNode, e *dbent.ProxyNode) {
	node.ID = e.ID
	node.IPAddress = e.IPAddress
	node.Country = e.Country
	node.CountryCode = e.CountryCode
	node.City = e.City
	node.ISP = e.Isp
	node.HTTPPort = e.HTTPPort
	node.VlessPort = e.VlessPort
	node.VlessNetwork = e.VlessNetwork
	node.VlessTLS = e.VlessTLS
	node.VlessSNI = e.VlessSni
	node.VlessWSPath = e.VlessWsPath
	node.Tags = e.Tags
	node.Status = e.Status
	node.Description = e.Description
	node.CreatedAt = e.CreatedAt
	node.UpdatedAt = e.UpdatedAt
}
