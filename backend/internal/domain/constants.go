package domain

// Status constants
const (
	StatusActive   = "active"
	StatusDisabled = "disabled"
	StatusError    = "error"
	StatusUnused   = "unused"
	StatusUsed     = "used"
	StatusExpired  = "expired"
)

// Role constants
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

// Redeem type constants
const (
	RedeemTypeBalance      = "balance"
	RedeemTypeConcurrency  = "concurrency"
	RedeemTypeSubscription = "subscription"
	RedeemTypeInvitation   = "invitation"
)

// PromoCode status constants
const (
	PromoCodeStatusActive   = "active"
	PromoCodeStatusDisabled = "disabled"
)

// Admin adjustment type constants
const (
	AdjustmentTypeAdminBalance     = "admin_balance"
	AdjustmentTypeAdminConcurrency = "admin_concurrency"
)

// Group subscription type constants
const (
	SubscriptionTypeStandard     = "standard"
	SubscriptionTypeSubscription = "subscription"
)

// Subscription status constants
const (
	SubscriptionStatusActive    = "active"
	SubscriptionStatusExpired   = "expired"
	SubscriptionStatusSuspended = "suspended"
)

// Payment order status constants
const (
	PaymentOrderStatusPending           = "pending"
	PaymentOrderStatusPaid              = "paid"
	PaymentOrderStatusRecharging        = "recharging"
	PaymentOrderStatusCompleted         = "completed"
	PaymentOrderStatusExpired           = "expired"
	PaymentOrderStatusCancelled         = "cancelled"
	PaymentOrderStatusFailed            = "failed"
	PaymentOrderStatusRefundRequested   = "refund_requested"
	PaymentOrderStatusRefunding         = "refunding"
	PaymentOrderStatusPartiallyRefunded = "partially_refunded"
	PaymentOrderStatusRefunded          = "refunded"
	PaymentOrderStatusRefundFailed      = "refund_failed"
)

// Payment order type constants
const (
	PaymentOrderTypeBalance      = "balance"
	PaymentOrderTypeSubscription = "subscription"
)

// Payment type constants
const (
	PaymentTypeAlipay       = "alipay"
	PaymentTypeAlipayDirect = "alipay_direct"
	PaymentTypeWxpay        = "wxpay"
	PaymentTypeWxpayDirect  = "wxpay_direct"
	PaymentTypeStripe       = "stripe"
)

// Proxy node status constants
const (
	ProxyNodeStatusAvailable = "available"
	ProxyNodeStatusRented    = "rented"
	ProxyNodeStatusOffline   = "offline"
)

// Proxy rental status constants
const (
	ProxyRentalStatusPendingPayment = "pending_payment"
	ProxyRentalStatusActive         = "active"
	ProxyRentalStatusExpired        = "expired"
	ProxyRentalStatusCancelled      = "cancelled"
)

// VLESS network type constants
const (
	VlessNetworkTCP  = "tcp"
	VlessNetworkWS   = "ws"
	VlessNetworkGRPC = "grpc"
)

// Payment order type for proxy rental
const (
	PaymentOrderTypeProxyRental = "proxy_rental"
)

// Payment provider key constants
const (
	PaymentProviderEasyPay = "easypay"
	PaymentProviderAlipay  = "alipay"
	PaymentProviderWxpay   = "wxpay"
	PaymentProviderStripe  = "stripe"
)
