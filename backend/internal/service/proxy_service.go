package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"math/big"
	"strings"
	"time"

	dbent "github.com/meitianwang/fast-frame/ent"
	"github.com/meitianwang/fast-frame/internal/domain"
	"github.com/meitianwang/fast-frame/internal/pkg/pagination"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ProxyService handles all proxy node and rental business logic.
type ProxyService struct {
	entClient   *dbent.Client
	nodeRepo    ProxyNodeRepository
	productRepo ProxyProductRepository
	rentalRepo  ProxyRentalRepository
	credRepo    ProxyCredentialRepository
	trafficRepo ProxyTrafficLogRepository
}

func NewProxyService(
	entClient *dbent.Client,
	nodeRepo ProxyNodeRepository,
	productRepo ProxyProductRepository,
	rentalRepo ProxyRentalRepository,
	credRepo ProxyCredentialRepository,
	trafficRepo ProxyTrafficLogRepository,
) *ProxyService {
	return &ProxyService{
		entClient:   entClient,
		nodeRepo:    nodeRepo,
		productRepo: productRepo,
		rentalRepo:  rentalRepo,
		credRepo:    credRepo,
		trafficRepo: trafficRepo,
	}
}

// --- Node management ---

func (s *ProxyService) CreateNode(ctx context.Context, node *ProxyNode) error {
	return s.nodeRepo.Create(ctx, node)
}

func (s *ProxyService) GetNode(ctx context.Context, id int64) (*ProxyNode, error) {
	return s.nodeRepo.GetByID(ctx, id)
}

func (s *ProxyService) UpdateNode(ctx context.Context, node *ProxyNode) error {
	return s.nodeRepo.Update(ctx, node)
}

func (s *ProxyService) DeleteNode(ctx context.Context, id int64) error {
	return s.nodeRepo.SoftDelete(ctx, id)
}

func (s *ProxyService) ListNodes(ctx context.Context, filter ProxyNodeFilter, params pagination.PaginationParams) ([]ProxyNode, *pagination.PaginationResult, error) {
	return s.nodeRepo.List(ctx, filter, params)
}

func (s *ProxyService) ListAvailableNodes(ctx context.Context, filter ProxyNodeFilter) ([]ProxyNode, error) {
	return s.nodeRepo.ListAvailable(ctx, filter)
}

// --- Product management ---

func (s *ProxyService) CreateProduct(ctx context.Context, p *ProxyProduct) error {
	return s.productRepo.Create(ctx, p)
}

func (s *ProxyService) GetProduct(ctx context.Context, id int64) (*ProxyProduct, error) {
	return s.productRepo.GetByID(ctx, id)
}

func (s *ProxyService) UpdateProduct(ctx context.Context, p *ProxyProduct) error {
	return s.productRepo.Update(ctx, p)
}

func (s *ProxyService) DeleteProduct(ctx context.Context, id int64) error {
	return s.productRepo.Delete(ctx, id)
}

func (s *ProxyService) ListProducts(ctx context.Context, activeOnly bool) ([]ProxyProduct, error) {
	return s.productRepo.List(ctx, activeOnly)
}

// --- Rental ---

// CreateRentalRequest holds inputs for initiating a rental.
type CreateRentalRequest struct {
	UserID    int64
	NodeID    int64
	ProductID int64
}

// CreateRental creates a pending rental within a transaction.
// It locks the node row to prevent concurrent purchases of the same IP,
// then returns the rental (status=pending_payment). Payment order creation
// is handled by the caller (payment service).
func (s *ProxyService) CreateRental(ctx context.Context, req CreateRentalRequest) (*ProxyRental, error) {
	product, err := s.productRepo.GetByID(ctx, req.ProductID)
	if err != nil {
		return nil, err
	}

	tx, err := s.entClient.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	txCtx := dbent.NewTxContext(ctx, tx)

	// Lock the node row and verify it's available.
	node, err := s.nodeRepo.LockForRental(txCtx, req.NodeID)
	if err != nil {
		return nil, err
	}
	if node.Status != domain.ProxyNodeStatusAvailable {
		return nil, ErrProxyNodeNotAvailable
	}

	trafficLimitBytes := int64(product.TrafficLimitGB) * 1024 * 1024 * 1024

	rental := &ProxyRental{
		UserID:            req.UserID,
		NodeID:            req.NodeID,
		ProductID:         req.ProductID,
		Status:            domain.ProxyRentalStatusPendingPayment,
		TrafficLimitBytes: trafficLimitBytes,
	}
	if err := s.rentalRepo.Create(txCtx, rental); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit: %w", err)
	}
	return rental, nil
}

// ActivateByPaymentOrder is called by the payment webhook after a successful payment.
// It activates the rental, marks the node as rented, and generates credentials.
func (s *ProxyService) ActivateByPaymentOrder(ctx context.Context, paymentOrderID int64) error {
	rental, err := s.rentalRepo.GetByPaymentOrderID(ctx, paymentOrderID)
	if err != nil {
		return err
	}
	if rental.Status != domain.ProxyRentalStatusPendingPayment {
		// Already activated or cancelled — idempotent.
		return nil
	}

	node, err := s.nodeRepo.GetByID(ctx, rental.NodeID)
	if err != nil {
		return err
	}
	product, err := s.productRepo.GetByID(ctx, rental.ProductID)
	if err != nil {
		return err
	}

	now := time.Now()
	expiresAt := now.AddDate(0, 0, product.DurationDays)

	tx, err := s.entClient.Tx(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	txCtx := dbent.NewTxContext(ctx, tx)

	rental.Status = domain.ProxyRentalStatusActive
	rental.StartedAt = &now
	rental.ExpiresAt = &expiresAt
	if err := s.rentalRepo.Update(txCtx, rental); err != nil {
		return err
	}

	if err := s.nodeRepo.SetStatus(txCtx, node.ID, domain.ProxyNodeStatusRented); err != nil {
		return err
	}

	cred, err := generateCredential(rental.ID, node, product)
	if err != nil {
		return fmt.Errorf("generate credential: %w", err)
	}
	if err := s.credRepo.Create(txCtx, cred); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	return nil
}

// SetRentalPaymentOrder links the payment order to the rental.
func (s *ProxyService) SetRentalPaymentOrder(ctx context.Context, rentalID, paymentOrderID int64) error {
	rental, err := s.rentalRepo.GetByID(ctx, rentalID)
	if err != nil {
		return err
	}
	rental.PaymentOrderID = &paymentOrderID
	return s.rentalRepo.Update(ctx, rental)
}

// CancelRental cancels a pending_payment rental and frees the node.
func (s *ProxyService) CancelRental(ctx context.Context, rentalID, userID int64) error {
	rental, err := s.rentalRepo.GetByID(ctx, rentalID)
	if err != nil {
		return err
	}
	if rental.UserID != userID {
		return ErrProxyRentalNotFound
	}
	if rental.Status != domain.ProxyRentalStatusPendingPayment {
		return ErrProxyRentalNotCancellable
	}
	rental.Status = domain.ProxyRentalStatusCancelled
	return s.rentalRepo.Update(ctx, rental)
}

// GetRentalWithCredential returns rental detail plus credential (only if active).
func (s *ProxyService) GetRentalWithCredential(ctx context.Context, rentalID, userID int64) (*ProxyRental, *ProxyCredential, error) {
	rental, err := s.rentalRepo.GetByID(ctx, rentalID)
	if err != nil {
		return nil, nil, err
	}
	if rental.UserID != userID {
		return nil, nil, ErrProxyRentalNotFound
	}

	var cred *ProxyCredential
	if rental.Status == domain.ProxyRentalStatusActive {
		cred, err = s.credRepo.GetByRentalID(ctx, rentalID)
		if err != nil {
			slog.Warn("credential not found for active rental", "rentalID", rentalID, "err", err)
		}
	}
	return rental, cred, nil
}

// ListUserRentals returns paginated rentals for a user.
func (s *ProxyService) ListUserRentals(ctx context.Context, userID int64, params pagination.PaginationParams) ([]ProxyRental, *pagination.PaginationResult, error) {
	return s.rentalRepo.ListByUserID(ctx, userID, params)
}

// ListRentals returns paginated rentals for admin view.
func (s *ProxyService) ListRentals(ctx context.Context, filter ProxyRentalFilter, params pagination.PaginationParams) ([]ProxyRental, *pagination.PaginationResult, error) {
	return s.rentalRepo.List(ctx, filter, params)
}

// AdminGetRental returns rental with credential for admin.
func (s *ProxyService) AdminGetRental(ctx context.Context, rentalID int64) (*ProxyRental, *ProxyCredential, error) {
	rental, err := s.rentalRepo.GetByID(ctx, rentalID)
	if err != nil {
		return nil, nil, err
	}
	var cred *ProxyCredential
	if rental.Status == domain.ProxyRentalStatusActive {
		cred, _ = s.credRepo.GetByRentalID(ctx, rentalID)
	}
	return rental, cred, nil
}

// UpdateTraffic adds delta bytes to the rental and auto-expires if limit exceeded.
func (s *ProxyService) UpdateTraffic(ctx context.Context, rentalID, operatorID, deltaBytes int64, note string) error {
	rental, err := s.rentalRepo.GetByID(ctx, rentalID)
	if err != nil {
		return err
	}

	log := &ProxyTrafficLog{
		RentalID:   rentalID,
		DeltaBytes: deltaBytes,
		OperatorID: operatorID,
		Note:       note,
	}
	if err := s.trafficRepo.Create(ctx, log); err != nil {
		return err
	}

	rental.TrafficUsedBytes += deltaBytes
	if rental.TrafficLimitBytes > 0 && rental.TrafficUsedBytes >= rental.TrafficLimitBytes {
		rental.Status = domain.ProxyRentalStatusExpired
		if err := s.rentalRepo.Update(ctx, rental); err != nil {
			return err
		}
		return s.nodeRepo.SetStatus(ctx, rental.NodeID, domain.ProxyNodeStatusAvailable)
	}

	return s.rentalRepo.Update(ctx, rental)
}

// ForceExpireRental expires a rental immediately (admin action).
func (s *ProxyService) ForceExpireRental(ctx context.Context, rentalID int64) error {
	rental, err := s.rentalRepo.GetByID(ctx, rentalID)
	if err != nil {
		return err
	}
	if rental.Status != domain.ProxyRentalStatusActive {
		return nil
	}
	rental.Status = domain.ProxyRentalStatusExpired
	if err := s.rentalRepo.Update(ctx, rental); err != nil {
		return err
	}
	return s.nodeRepo.SetStatus(ctx, rental.NodeID, domain.ProxyNodeStatusAvailable)
}

// ExpireOverdueRentals is called by the background scheduler.
func (s *ProxyService) ExpireOverdueRentals(ctx context.Context) {
	rentals, err := s.rentalRepo.ListExpiredActive(ctx, time.Now(), 200)
	if err != nil {
		slog.Error("proxy: list expired rentals failed", "err", err)
		return
	}
	for _, rental := range rentals {
		rental := rental
		rental.Status = domain.ProxyRentalStatusExpired
		if err := s.rentalRepo.Update(ctx, &rental); err != nil {
			slog.Error("proxy: expire rental failed", "rentalID", rental.ID, "err", err)
			continue
		}
		if err := s.nodeRepo.SetStatus(ctx, rental.NodeID, domain.ProxyNodeStatusAvailable); err != nil {
			slog.Error("proxy: free node after expiry failed", "nodeID", rental.NodeID, "err", err)
		}
	}
	if len(rentals) > 0 {
		slog.Info("proxy: expired overdue rentals", "count", len(rentals))
	}
}

// GetTrafficLogs returns traffic update history for a rental.
func (s *ProxyService) GetTrafficLogs(ctx context.Context, rentalID int64) ([]ProxyTrafficLog, error) {
	return s.trafficRepo.ListByRentalID(ctx, rentalID)
}

// --- Credential generation ---

const credChars = "abcdefghjkmnpqrstuvwxyzABCDEFGHJKMNPQRSTUVWXYZ23456789"
const passChars = "abcdefghjkmnpqrstuvwxyzABCDEFGHJKMNPQRSTUVWXYZ23456789!@#$%^&*"

func randomString(chars string, n int) string {
	b := make([]byte, n)
	for i := range b {
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		b[i] = chars[idx.Int64()]
	}
	return string(b)
}

func generateCredential(rentalID int64, node *ProxyNode, product *ProxyProduct) (*ProxyCredential, error) {
	vlessUUID := uuid.New().String()
	username := randomString(credChars, 8)
	password := randomString(passChars, 16)

	label := node.City
	if label == "" {
		label = node.Country
	}
	if label == "" {
		label = node.IPAddress
	}

	vlessLink := buildVlessLink(vlessUUID, node, label+"-家庭IP")

	return &ProxyCredential{
		RentalID:     rentalID,
		HTTPUsername: username,
		HTTPPassword: password,
		VlessUUID:    vlessUUID,
		VlessLink:    vlessLink,
	}, nil
}

// buildVlessLink constructs the vless:// URI from node configuration.
func buildVlessLink(uuid string, node *ProxyNode, name string) string {
	security := "none"
	if node.VlessTLS {
		security = "tls"
	}

	params := []string{
		"encryption=none",
		"security=" + security,
		"type=" + node.VlessNetwork,
	}

	if node.VlessTLS && node.VlessSNI != "" {
		params = append(params, "sni="+node.VlessSNI)
		params = append(params, "host="+node.VlessSNI)
	}
	if node.VlessNetwork == domain.VlessNetworkWS && node.VlessWSPath != "" {
		params = append(params, "path="+node.VlessWSPath)
	}

	return fmt.Sprintf("vless://%s@%s:%d?%s#%s",
		uuid, node.IPAddress, node.VlessPort,
		strings.Join(params, "&"),
		name,
	)
}

// PriceForRental returns the price for a given product (convenience for payment order creation).
func (s *ProxyService) PriceForProduct(ctx context.Context, productID int64) (decimal.Decimal, error) {
	p, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return decimal.Zero, err
	}
	return p.Price, nil
}
