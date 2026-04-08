package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/ent/paymentorder"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/shopspring/decimal"
)

type paymentOrderRepository struct {
	client *dbent.Client
	sql    *sql.DB
}

func NewPaymentOrderRepository(client *dbent.Client, sqlDB *sql.DB) service.PaymentOrderRepository {
	return &paymentOrderRepository{client: client, sql: sqlDB}
}

func (r *paymentOrderRepository) Create(ctx context.Context, order *service.PaymentOrder) error {
	client := clientFromContext(ctx, r.client)
	builder := client.PaymentOrder.Create().
		SetUserID(order.UserID).
		SetNillableUserEmail(order.UserEmail).
		SetNillableUserName(order.UserName).
		SetNillableUserNotes(order.UserNotes).
		SetAmount(order.Amount).
		SetNillablePayAmount(order.PayAmount).
		SetNillableFeeRate(order.FeeRate).
		SetRechargeCode(order.RechargeCode).
		SetStatus(order.Status).
		SetPaymentType(order.PaymentType).
		SetNillablePaymentTradeNo(order.PaymentTradeNo).
		SetNillablePayURL(order.PayURL).
		SetNillableQrCode(order.QrCode).
		SetExpiresAt(order.ExpiresAt).
		SetNillableClientIP(order.ClientIP).
		SetNillableSrcHost(order.SrcHost).
		SetNillableSrcURL(order.SrcURL).
		SetOrderType(order.OrderType).
		SetNillablePlanID(order.PlanID).
		SetNillableSubscriptionGroupID(order.SubscriptionGroupID).
		SetNillableSubscriptionDays(order.SubscriptionDays).
		SetNillableProviderInstanceID(order.ProviderInstanceID)

	created, err := builder.Save(ctx)
	if err != nil {
		return translatePersistenceError(err, nil, nil)
	}
	applyPaymentOrderEntity(order, created)
	return nil
}

func (r *paymentOrderRepository) Delete(ctx context.Context, id int64) error {
	client := clientFromContext(ctx, r.client)
	return client.PaymentOrder.DeleteOneID(id).Exec(ctx)
}

func (r *paymentOrderRepository) GetByID(ctx context.Context, id int64) (*service.PaymentOrder, error) {
	client := clientFromContext(ctx, r.client)
	m, err := client.PaymentOrder.Get(ctx, id)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrPaymentOrderNotFound, nil)
	}
	return paymentOrderToService(m), nil
}

func (r *paymentOrderRepository) GetByRechargeCode(ctx context.Context, code string) (*service.PaymentOrder, error) {
	client := clientFromContext(ctx, r.client)
	m, err := client.PaymentOrder.Query().
		Where(paymentorder.RechargeCodeEQ(code)).
		Only(ctx)
	if err != nil {
		return nil, translatePersistenceError(err, service.ErrPaymentOrderNotFound, nil)
	}
	return paymentOrderToService(m), nil
}

func (r *paymentOrderRepository) Update(ctx context.Context, order *service.PaymentOrder) error {
	client := clientFromContext(ctx, r.client)
	builder := client.PaymentOrder.UpdateOneID(order.ID).
		SetStatus(order.Status).
		SetNillablePaymentTradeNo(order.PaymentTradeNo).
		SetNillablePayURL(order.PayURL).
		SetNillableQrCode(order.QrCode).
		SetNillableRefundAmount(order.RefundAmount).
		SetNillableRefundReason(order.RefundReason).
		SetNillableRefundAt(order.RefundAt).
		SetForceRefund(order.ForceRefund).
		SetNillableRefundRequestedAt(order.RefundRequestedAt).
		SetNillableRefundRequestReason(order.RefundRequestReason).
		SetNillableRefundRequestedBy(order.RefundRequestedBy).
		SetNillablePaidAt(order.PaidAt).
		SetNillableCompletedAt(order.CompletedAt).
		SetNillableFailedAt(order.FailedAt).
		SetNillableFailedReason(order.FailedReason).
		SetNillableProviderInstanceID(order.ProviderInstanceID)

	_, err := builder.Save(ctx)
	return translatePersistenceError(err, service.ErrPaymentOrderNotFound, nil)
}

func (r *paymentOrderRepository) UpdateRechargeCode(ctx context.Context, id int64, code string) error {
	client := clientFromContext(ctx, r.client)
	_, err := client.PaymentOrder.UpdateOneID(id).
		SetRechargeCode(code).
		Save(ctx)
	return translatePersistenceError(err, service.ErrPaymentOrderNotFound, nil)
}

func (r *paymentOrderRepository) UpdateStatusCAS(ctx context.Context, id int64, fromStatus, toStatus string) (bool, error) {
	client := clientFromContext(ctx, r.client)
	n, err := client.PaymentOrder.Update().
		Where(paymentorder.IDEQ(id), paymentorder.StatusEQ(fromStatus)).
		SetStatus(toStatus).
		Save(ctx)
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (r *paymentOrderRepository) ConfirmPaidCAS(ctx context.Context, id int64, tradeNo string, paidAmount decimal.Decimal, paidAt time.Time, graceDeadline time.Time) (bool, error) {
	client := clientFromContext(ctx, r.client)
	n, err := client.PaymentOrder.Update().
		Where(
			paymentorder.IDEQ(id),
			paymentorder.Or(
				paymentorder.StatusEQ("pending"),
				paymentorder.And(
					paymentorder.StatusEQ("expired"),
					paymentorder.UpdatedAtGTE(graceDeadline),
				),
			),
		).
		SetStatus("paid").
		SetPaymentTradeNo(tradeNo).
		SetPaidAt(paidAt).
		ClearFailedAt().
		ClearFailedReason().
		Save(ctx)
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

func (r *paymentOrderRepository) UpdatePaymentResult(ctx context.Context, id int64, tradeNo string, payURL, qrCode *string, instanceID *int64) error {
	client := clientFromContext(ctx, r.client)
	builder := client.PaymentOrder.UpdateOneID(id).
		SetPaymentTradeNo(tradeNo)
	if payURL != nil {
		builder.SetPayURL(*payURL)
	}
	if qrCode != nil {
		builder.SetQrCode(*qrCode)
	}
	if instanceID != nil {
		builder.SetProviderInstanceID(*instanceID)
	}
	_, err := builder.Save(ctx)
	return translatePersistenceError(err, service.ErrPaymentOrderNotFound, nil)
}

func (r *paymentOrderRepository) CountPendingByUserID(ctx context.Context, userID int64) (int, error) {
	client := clientFromContext(ctx, r.client)
	return client.PaymentOrder.Query().
		Where(paymentorder.UserIDEQ(userID), paymentorder.StatusEQ("pending")).
		Count(ctx)
}

// SumDailyPaidByUserID sums base recharge amounts (amount, not pay_amount) for user daily limit checks.
// User-facing limits are denominated in recharge value, not the fee-inclusive amount charged at the gateway.
func (r *paymentOrderRepository) SumDailyPaidByUserID(ctx context.Context, userID int64, bizDayStart time.Time) (decimal.Decimal, error) {
	return r.sumAmount(ctx,
		"SELECT COALESCE(SUM(amount), 0) FROM payment_orders WHERE user_id = $1 AND paid_at >= $2 AND status IN ('paid', 'recharging', 'completed') ",
		userID, bizDayStart,
	)
}

// SumDailyPaidByPaymentType sums base recharge amounts for global per-method daily limit checks.
func (r *paymentOrderRepository) SumDailyPaidByPaymentType(ctx context.Context, paymentType string, bizDayStart time.Time) (decimal.Decimal, error) {
	return r.sumAmount(ctx,
		"SELECT COALESCE(SUM(amount), 0) FROM payment_orders WHERE payment_type = $1 AND paid_at >= $2 AND status IN ('paid', 'recharging', 'completed') ",
		paymentType, bizDayStart,
	)
}

// SumDailyPaidByInstanceID sums fee-inclusive amounts (pay_amount) for provider instance daily limit checks.
// Instance limits track actual money flowing through the payment gateway, which includes fees.
func (r *paymentOrderRepository) SumDailyPaidByInstanceID(ctx context.Context, instanceID int64, bizDayStart time.Time) (decimal.Decimal, error) {
	return r.sumAmount(ctx,
		"SELECT COALESCE(SUM(pay_amount), 0) FROM payment_orders WHERE provider_instance_id = $1 AND paid_at >= $2 AND status IN ('paid', 'recharging', 'completed') ",
		instanceID, bizDayStart,
	)
}

// SumDailyPaidByInstanceIDs is the batch version of SumDailyPaidByInstanceID (uses pay_amount).
func (r *paymentOrderRepository) SumDailyPaidByInstanceIDs(ctx context.Context, instanceIDs []int64, bizDayStart time.Time) (map[int64]decimal.Decimal, error) {
	result := make(map[int64]decimal.Decimal, len(instanceIDs))
	if len(instanceIDs) == 0 {
		return result, nil
	}

	// Build placeholders: $1, $2, ..., $N for instance IDs; $N+1 for bizDayStart
	placeholders := make([]string, len(instanceIDs))
	args := make([]any, 0, len(instanceIDs)+1)
	for i, id := range instanceIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args = append(args, id)
	}
	args = append(args, bizDayStart)
	tsPlaceholder := fmt.Sprintf("$%d", len(instanceIDs)+1)

	query := fmt.Sprintf(
		"SELECT provider_instance_id, COALESCE(SUM(pay_amount), 0) FROM payment_orders WHERE provider_instance_id IN (%s) AND paid_at >= %s AND status IN ('paid', 'recharging', 'completed') GROUP BY provider_instance_id",
		strings.Join(placeholders, ","), tsPlaceholder,
	)

	rows, err := r.sql.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var instID int64
		var s string
		if err := rows.Scan(&instID, &s); err != nil {
			return nil, err
		}
		d, _ := decimal.NewFromString(s)
		result[instID] = d
	}
	return result, rows.Err()
}

func (r *paymentOrderRepository) ListExpiredPending(ctx context.Context, now time.Time, limit int) ([]service.PaymentOrder, error) {
	client := clientFromContext(ctx, r.client)
	ms, err := client.PaymentOrder.Query().
		Where(
			paymentorder.StatusEQ("pending"),
			paymentorder.ExpiresAtLT(now),
		).
		Limit(limit).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return paymentOrdersToService(ms), nil
}

func (r *paymentOrderRepository) List(ctx context.Context, params pagination.PaginationParams) ([]service.PaymentOrder, *pagination.PaginationResult, error) {
	client := clientFromContext(ctx, r.client)
	query := client.PaymentOrder.Query().Order(dbent.Desc(paymentorder.FieldID))

	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, nil, err
	}

	ms, err := query.
		Offset((params.Page - 1) * params.PageSize).
		Limit(params.PageSize).
		All(ctx)
	if err != nil {
		return nil, nil, err
	}

	return paymentOrdersToService(ms), &pagination.PaginationResult{
		Total:    int64(total),
		Page:     params.Page,
		PageSize: params.PageSize,
	}, nil
}

func (r *paymentOrderRepository) ListByUserID(ctx context.Context, userID int64, params pagination.PaginationParams) ([]service.PaymentOrder, *pagination.PaginationResult, error) {
	client := clientFromContext(ctx, r.client)
	query := client.PaymentOrder.Query().
		Where(paymentorder.UserIDEQ(userID)).
		Order(dbent.Desc(paymentorder.FieldID))

	total, err := query.Clone().Count(ctx)
	if err != nil {
		return nil, nil, err
	}

	ms, err := query.
		Offset((params.Page - 1) * params.PageSize).
		Limit(params.PageSize).
		All(ctx)
	if err != nil {
		return nil, nil, err
	}

	return paymentOrdersToService(ms), &pagination.PaginationResult{
		Total:    int64(total),
		Page:     params.Page,
		PageSize: params.PageSize,
	}, nil
}

// --- helpers ---

func (r *paymentOrderRepository) sumAmount(ctx context.Context, query string, args ...any) (decimal.Decimal, error) {
	row := r.sql.QueryRowContext(ctx, query, args...)
	var s string
	if err := row.Scan(&s); err != nil {
		return decimal.Zero, err
	}
	d, err := decimal.NewFromString(s)
	if err != nil {
		return decimal.Zero, nil
	}
	return d, nil
}

func paymentOrderToService(m *dbent.PaymentOrder) *service.PaymentOrder {
	if m == nil {
		return nil
	}
	return &service.PaymentOrder{
		ID:                  m.ID,
		UserID:              m.UserID,
		UserEmail:           m.UserEmail,
		UserName:            m.UserName,
		UserNotes:           m.UserNotes,
		Amount:              m.Amount,
		PayAmount:           m.PayAmount,
		FeeRate:             m.FeeRate,
		RechargeCode:        m.RechargeCode,
		Status:              m.Status,
		PaymentType:         m.PaymentType,
		PaymentTradeNo:      m.PaymentTradeNo,
		PayURL:              m.PayURL,
		QrCode:              m.QrCode,
		QrCodeImg:           m.QrCodeImg,
		RefundAmount:        m.RefundAmount,
		RefundReason:        m.RefundReason,
		RefundAt:            m.RefundAt,
		ForceRefund:         m.ForceRefund,
		RefundRequestedAt:   m.RefundRequestedAt,
		RefundRequestReason: m.RefundRequestReason,
		RefundRequestedBy:   m.RefundRequestedBy,
		ExpiresAt:           m.ExpiresAt,
		PaidAt:              m.PaidAt,
		CompletedAt:         m.CompletedAt,
		FailedAt:            m.FailedAt,
		FailedReason:        m.FailedReason,
		ClientIP:            m.ClientIP,
		SrcHost:             m.SrcHost,
		SrcURL:              m.SrcURL,
		OrderType:           m.OrderType,
		PlanID:              m.PlanID,
		SubscriptionGroupID: m.SubscriptionGroupID,
		SubscriptionDays:    m.SubscriptionDays,
		ProviderInstanceID:  m.ProviderInstanceID,
		CreatedAt:           m.CreatedAt,
		UpdatedAt:           m.UpdatedAt,
	}
}

func applyPaymentOrderEntity(dst *service.PaymentOrder, src *dbent.PaymentOrder) {
	dst.ID = src.ID
	dst.CreatedAt = src.CreatedAt
	dst.UpdatedAt = src.UpdatedAt
}

func paymentOrdersToService(ms []*dbent.PaymentOrder) []service.PaymentOrder {
	result := make([]service.PaymentOrder, len(ms))
	for i, m := range ms {
		result[i] = *paymentOrderToService(m)
	}
	return result
}
