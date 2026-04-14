package service

import (
	"context"
	"time"

	infraerrors "github.com/meitianwang/fast-frame/internal/pkg/errors"
	"github.com/meitianwang/fast-frame/internal/pkg/pagination"
	"github.com/shopspring/decimal"
)

// --- Sentinel errors ---

var (
	ErrProxyNodeNotFound     = infraerrors.NotFound("PROXY_NODE_NOT_FOUND", "proxy node not found")
	ErrProxyProductNotFound  = infraerrors.NotFound("PROXY_PRODUCT_NOT_FOUND", "proxy product not found")
	ErrProxyRentalNotFound   = infraerrors.NotFound("PROXY_RENTAL_NOT_FOUND", "proxy rental not found")
	ErrProxyNodeNotAvailable    = infraerrors.Conflict("PROXY_NODE_NOT_AVAILABLE", "proxy node is not available for rental")
	ErrProxyRentalNotCancellable = infraerrors.Conflict("RENTAL_NOT_CANCELLABLE", "rental cannot be cancelled in current status")
)

// --- Domain types ---

// ProxyNode represents a VPS-based residential proxy node.
type ProxyNode struct {
	ID           int64
	IPAddress    string
	Country      string
	CountryCode  string
	City         string
	ISP          string
	HTTPPort     int
	VlessPort    int
	VlessNetwork string
	VlessTLS     bool
	VlessSNI     string
	VlessWSPath  string
	Tags         []string
	Status       string
	Description  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// ProxyProduct represents a rental plan.
type ProxyProduct struct {
	ID             int64
	Name           string
	Description    string
	DurationDays   int
	TrafficLimitGB int
	Price          decimal.Decimal
	SortOrder      int
	IsActive       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// ProxyRental represents a user's active or historical IP rental.
type ProxyRental struct {
	ID                 int64
	UserID             int64
	NodeID             int64
	ProductID          int64
	PaymentOrderID     *int64
	Status             string
	StartedAt          *time.Time
	ExpiresAt          *time.Time
	TrafficUsedBytes   int64
	TrafficLimitBytes  int64
	CreatedAt          time.Time
	UpdatedAt          time.Time

	// Eagerly loaded associations (optional)
	Node    *ProxyNode
	Product *ProxyProduct
}

// ProxyCredential holds generated access credentials for a rental.
type ProxyCredential struct {
	ID           int64
	RentalID     int64
	HTTPUsername string
	HTTPPassword string
	VlessUUID    string
	VlessLink    string
	CreatedAt    time.Time
}

// ProxyTrafficLog records an admin traffic usage update.
type ProxyTrafficLog struct {
	ID          int64
	RentalID    int64
	DeltaBytes  int64
	OperatorID  int64
	Note        string
	CreatedAt   time.Time
}

// --- Repository interfaces ---

type ProxyNodeRepository interface {
	Create(ctx context.Context, node *ProxyNode) error
	GetByID(ctx context.Context, id int64) (*ProxyNode, error)
	Update(ctx context.Context, node *ProxyNode) error
	SoftDelete(ctx context.Context, id int64) error
	List(ctx context.Context, filter ProxyNodeFilter, params pagination.PaginationParams) ([]ProxyNode, *pagination.PaginationResult, error)
	ListAvailable(ctx context.Context, filter ProxyNodeFilter) ([]ProxyNode, error)
	// LockForRental locks the node row in a transaction and returns the current status.
	LockForRental(ctx context.Context, id int64) (*ProxyNode, error)
	SetStatus(ctx context.Context, id int64, status string) error
}

type ProxyProductRepository interface {
	Create(ctx context.Context, p *ProxyProduct) error
	GetByID(ctx context.Context, id int64) (*ProxyProduct, error)
	Update(ctx context.Context, p *ProxyProduct) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, activeOnly bool) ([]ProxyProduct, error)
}

type ProxyRentalRepository interface {
	Create(ctx context.Context, r *ProxyRental) error
	GetByID(ctx context.Context, id int64) (*ProxyRental, error)
	GetByPaymentOrderID(ctx context.Context, orderID int64) (*ProxyRental, error)
	Update(ctx context.Context, r *ProxyRental) error
	List(ctx context.Context, filter ProxyRentalFilter, params pagination.PaginationParams) ([]ProxyRental, *pagination.PaginationResult, error)
	ListByUserID(ctx context.Context, userID int64, params pagination.PaginationParams) ([]ProxyRental, *pagination.PaginationResult, error)
	// ListExpiredActive returns active rentals whose expires_at is past now.
	ListExpiredActive(ctx context.Context, now time.Time, limit int) ([]ProxyRental, error)
}

type ProxyCredentialRepository interface {
	Create(ctx context.Context, cred *ProxyCredential) error
	GetByRentalID(ctx context.Context, rentalID int64) (*ProxyCredential, error)
}

type ProxyTrafficLogRepository interface {
	Create(ctx context.Context, log *ProxyTrafficLog) error
	ListByRentalID(ctx context.Context, rentalID int64) ([]ProxyTrafficLog, error)
}

// --- Filter types ---

type ProxyNodeFilter struct {
	CountryCode string
	Tag         string
	Status      string
}

type ProxyRentalFilter struct {
	UserID  *int64
	NodeID  *int64
	Status  string
}
