package repository

import (
	"context"

	dbent "github.com/meitianwang/fast-frame/ent"
	"github.com/meitianwang/fast-frame/ent/proxycredential"
	"github.com/meitianwang/fast-frame/ent/proxytrafficlog"
	"github.com/meitianwang/fast-frame/internal/service"
)

type proxyCredentialRepository struct {
	client *dbent.Client
}

func NewProxyCredentialRepository(client *dbent.Client) service.ProxyCredentialRepository {
	return &proxyCredentialRepository{client: client}
}

func (r *proxyCredentialRepository) Create(ctx context.Context, cred *service.ProxyCredential) error {
	client := clientFromContext(ctx, r.client)
	created, err := client.ProxyCredential.Create().
		SetRentalID(cred.RentalID).
		SetHTTPUsername(cred.HTTPUsername).
		SetHTTPPassword(cred.HTTPPassword).
		SetVlessUUID(cred.VlessUUID).
		SetVlessLink(cred.VlessLink).
		Save(ctx)
	if err != nil {
		return translatePersistenceError(err, nil, nil)
	}
	cred.ID = created.ID
	cred.CreatedAt = created.CreatedAt
	return nil
}

func (r *proxyCredentialRepository) GetByRentalID(ctx context.Context, rentalID int64) (*service.ProxyCredential, error) {
	client := clientFromContext(ctx, r.client)
	e, err := client.ProxyCredential.Query().
		Where(proxycredential.RentalID(rentalID)).
		Only(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrProxyRentalNotFound, nil)
	}
	return &service.ProxyCredential{
		ID:           e.ID,
		RentalID:     e.RentalID,
		HTTPUsername: e.HTTPUsername,
		HTTPPassword: e.HTTPPassword,
		VlessUUID:    e.VlessUUID,
		VlessLink:    e.VlessLink,
		CreatedAt:    e.CreatedAt,
	}, nil
}

type proxyTrafficLogRepository struct {
	client *dbent.Client
}

func NewProxyTrafficLogRepository(client *dbent.Client) service.ProxyTrafficLogRepository {
	return &proxyTrafficLogRepository{client: client}
}

func (r *proxyTrafficLogRepository) Create(ctx context.Context, log *service.ProxyTrafficLog) error {
	client := clientFromContext(ctx, r.client)
	created, err := client.ProxyTrafficLog.Create().
		SetRentalID(log.RentalID).
		SetDeltaBytes(log.DeltaBytes).
		SetOperatorID(log.OperatorID).
		SetNote(log.Note).
		Save(ctx)
	if err != nil {
		return translatePersistenceError(err, nil, nil)
	}
	log.ID = created.ID
	log.CreatedAt = created.CreatedAt
	return nil
}

func (r *proxyTrafficLogRepository) ListByRentalID(ctx context.Context, rentalID int64) ([]service.ProxyTrafficLog, error) {
	client := clientFromContext(ctx, r.client)
	entities, err := client.ProxyTrafficLog.Query().
		Where(proxytrafficlog.RentalID(rentalID)).
		Order(dbent.Desc(proxytrafficlog.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	logs := make([]service.ProxyTrafficLog, len(entities))
	for i, e := range entities {
		logs[i] = service.ProxyTrafficLog{
			ID:         e.ID,
			RentalID:   e.RentalID,
			DeltaBytes: e.DeltaBytes,
			OperatorID: e.OperatorID,
			Note:       e.Note,
			CreatedAt:  e.CreatedAt,
		}
	}
	return logs, nil
}
