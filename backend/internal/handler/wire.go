package handler

import (
	"github.com/meitianwang/fast-frame/internal/handler/admin"
	"github.com/meitianwang/fast-frame/internal/service"

	"github.com/google/wire"
)

// ProvideAdminHandlers creates the AdminHandlers struct
func ProvideAdminHandlers(
	proxyAdminHandler *admin.ProxyAdminHandler,
	dashboardHandler *admin.DashboardHandler,
	userHandler *admin.UserHandler,
	announcementHandler *admin.AnnouncementHandler,
	backupHandler *admin.BackupHandler,
	redeemHandler *admin.RedeemHandler,
	promoHandler *admin.PromoHandler,
	settingHandler *admin.SettingHandler,
	systemHandler *admin.SystemHandler,
	subscriptionHandler *admin.SubscriptionHandler,
	userAttributeHandler *admin.UserAttributeHandler,
	paymentOrderHandler *admin.PaymentOrderHandler,
	paymentRefundHandler *admin.PaymentRefundHandler,
	paymentConfigHandler *admin.PaymentConfigHandler,
	paymentProviderInstanceHandler *admin.PaymentProviderInstanceHandler,
	paymentChannelHandler *admin.PaymentChannelHandler,
	paymentSubscriptionPlanHandler *admin.PaymentSubscriptionPlanHandler,
	paymentDashboardHandler *admin.PaymentDashboardHandler,
) *AdminHandlers {
	return &AdminHandlers{
		Proxy:                 proxyAdminHandler,
		Dashboard:             dashboardHandler,
		User:                  userHandler,
		Announcement:          announcementHandler,
		Backup:                backupHandler,
		Redeem:                redeemHandler,
		Promo:                 promoHandler,
		Setting:               settingHandler,
		System:                systemHandler,
		Subscription:          subscriptionHandler,
		UserAttribute:         userAttributeHandler,
		PaymentOrder:            paymentOrderHandler,
		PaymentRefund:           paymentRefundHandler,
		PaymentConfig:           paymentConfigHandler,
		PaymentProviderInstance: paymentProviderInstanceHandler,
		PaymentChannel:          paymentChannelHandler,
		PaymentSubscriptionPlan: paymentSubscriptionPlanHandler,
		PaymentDashboard:        paymentDashboardHandler,
	}
}

// ProvideSystemHandler creates admin.SystemHandler with UpdateService
func ProvideSystemHandler(updateService *service.UpdateService, lockService *service.SystemOperationLockService) *admin.SystemHandler {
	return admin.NewSystemHandler(updateService, lockService)
}

// ProvideSettingHandler creates SettingHandler with version from BuildInfo
func ProvideSettingHandler(settingService *service.SettingService, buildInfo BuildInfo) *SettingHandler {
	return NewSettingHandler(settingService, buildInfo.Version)
}

// ProvideHandlers creates the Handlers struct
func ProvideHandlers(
	authHandler *AuthHandler,
	userHandler *UserHandler,
	redeemHandler *RedeemHandler,
	subscriptionHandler *SubscriptionHandler,
	announcementHandler *AnnouncementHandler,
	adminHandlers *AdminHandlers,
	settingHandler *SettingHandler,
	totpHandler *TotpHandler,
	paymentHandler *PaymentHandler,
	paymentWebhookHandler *PaymentWebhookHandler,
	proxyHandler *ProxyHandler,
) *Handlers {
	return &Handlers{
		Auth:           authHandler,
		User:           userHandler,
		Redeem:         redeemHandler,
		Subscription:   subscriptionHandler,
		Announcement:   announcementHandler,
		Admin:          adminHandlers,
		Setting:        settingHandler,
		Totp:           totpHandler,
		Payment:        paymentHandler,
		PaymentWebhook: paymentWebhookHandler,
		Proxy:          proxyHandler,
	}
}

// ProviderSet is the Wire provider set for all handlers
var ProviderSet = wire.NewSet(
	// Top-level handlers
	NewAuthHandler,
	NewUserHandler,
	NewRedeemHandler,
	NewSubscriptionHandler,
	NewAnnouncementHandler,
	NewTotpHandler,
	ProvideSettingHandler,

	// Admin handlers
	admin.NewDashboardHandler,
	admin.NewUserHandler,
	admin.NewAnnouncementHandler,
	admin.NewBackupHandler,
	admin.NewRedeemHandler,
	admin.NewPromoHandler,
	admin.NewSettingHandler,
	ProvideSystemHandler,
	admin.NewSubscriptionHandler,
	admin.NewUserAttributeHandler,

	// Proxy handlers
	NewProxyHandler,
	admin.NewProxyAdminHandler,

	// Payment handlers
	NewPaymentHandler,
	NewPaymentWebhookHandler,
	admin.NewPaymentOrderHandler,
	admin.NewPaymentRefundHandler,
	admin.NewPaymentConfigHandler,
	admin.NewPaymentProviderInstanceHandler,
	admin.NewPaymentChannelHandler,
	admin.NewPaymentSubscriptionPlanHandler,
	admin.NewPaymentDashboardHandler,

	// AdminHandlers and Handlers constructors
	ProvideAdminHandlers,
	ProvideHandlers,
)
