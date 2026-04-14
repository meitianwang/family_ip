package handler

import (
	"github.com/meitianwang/fast-frame/internal/handler/admin"
)

// AdminHandlers contains all admin-related HTTP handlers
type AdminHandlers struct {
	Proxy                 *admin.ProxyAdminHandler
	Dashboard             *admin.DashboardHandler
	User                  *admin.UserHandler
	Announcement          *admin.AnnouncementHandler
	Backup                *admin.BackupHandler
	Redeem                *admin.RedeemHandler
	Promo                 *admin.PromoHandler
	Setting               *admin.SettingHandler
	System                *admin.SystemHandler
	Subscription          *admin.SubscriptionHandler
	UserAttribute         *admin.UserAttributeHandler
	PaymentOrder            *admin.PaymentOrderHandler
	PaymentRefund           *admin.PaymentRefundHandler
	PaymentConfig           *admin.PaymentConfigHandler
	PaymentProviderInstance *admin.PaymentProviderInstanceHandler
	PaymentChannel          *admin.PaymentChannelHandler
	PaymentSubscriptionPlan *admin.PaymentSubscriptionPlanHandler
	PaymentDashboard        *admin.PaymentDashboardHandler
}

// Handlers contains all HTTP handlers
type Handlers struct {
	Auth           *AuthHandler
	User           *UserHandler
	Redeem         *RedeemHandler
	Subscription   *SubscriptionHandler
	Announcement   *AnnouncementHandler
	Admin          *AdminHandlers
	Setting        *SettingHandler
	Totp           *TotpHandler
	Payment        *PaymentHandler
	PaymentWebhook *PaymentWebhookHandler
	Proxy          *ProxyHandler
}

// BuildInfo contains build-time information
type BuildInfo struct {
	Version   string
	BuildType string // "source" for manual builds, "release" for CI builds
}
