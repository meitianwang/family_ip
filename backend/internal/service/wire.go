package service

import (
	"context"
	"time"

	dbent "github.com/meitianwang/fast-frame/ent"
	"github.com/meitianwang/fast-frame/internal/config"
	"github.com/google/wire"
)

// BuildInfo contains build information
type BuildInfo struct {
	Version   string
	BuildType string
}

// ProvideUpdateService creates UpdateService with BuildInfo
func ProvideUpdateService(cache UpdateCache, githubClient GitHubReleaseClient, buildInfo BuildInfo) *UpdateService {
	return NewUpdateService(cache, githubClient, buildInfo.Version, buildInfo.BuildType)
}

// ProvideEmailQueueService creates EmailQueueService with default worker count
func ProvideEmailQueueService(emailService *EmailService) *EmailQueueService {
	return NewEmailQueueService(emailService, 3)
}

// ProvideSubscriptionExpiryService creates and starts SubscriptionExpiryService.
func ProvideSubscriptionExpiryService(userSubRepo UserSubscriptionRepository) *SubscriptionExpiryService {
	svc := NewSubscriptionExpiryService(userSubRepo, time.Minute)
	svc.Start()
	return svc
}

func ProvideIdempotencyCoordinator(repo IdempotencyRepository, cfg *config.Config) *IdempotencyCoordinator {
	coordinator := NewIdempotencyCoordinator(repo, buildIdempotencyConfig(cfg))
	SetDefaultIdempotencyCoordinator(coordinator)
	return coordinator
}

func ProvideSystemOperationLockService(repo IdempotencyRepository, cfg *config.Config) *SystemOperationLockService {
	return NewSystemOperationLockService(repo, buildIdempotencyConfig(cfg))
}

func buildIdempotencyConfig(cfg *config.Config) IdempotencyConfig {
	idempotencyCfg := DefaultIdempotencyConfig()
	if cfg != nil {
		if cfg.Idempotency.DefaultTTLSeconds > 0 {
			idempotencyCfg.DefaultTTL = time.Duration(cfg.Idempotency.DefaultTTLSeconds) * time.Second
		}
		if cfg.Idempotency.SystemOperationTTLSeconds > 0 {
			idempotencyCfg.SystemOperationTTL = time.Duration(cfg.Idempotency.SystemOperationTTLSeconds) * time.Second
		}
		if cfg.Idempotency.ProcessingTimeoutSeconds > 0 {
			idempotencyCfg.ProcessingTimeout = time.Duration(cfg.Idempotency.ProcessingTimeoutSeconds) * time.Second
		}
		if cfg.Idempotency.FailedRetryBackoffSeconds > 0 {
			idempotencyCfg.FailedRetryBackoff = time.Duration(cfg.Idempotency.FailedRetryBackoffSeconds) * time.Second
		}
		if cfg.Idempotency.MaxStoredResponseLen > 0 {
			idempotencyCfg.MaxStoredResponseLen = cfg.Idempotency.MaxStoredResponseLen
		}
		idempotencyCfg.ObserveOnly = cfg.Idempotency.ObserveOnly
	}
	return idempotencyCfg
}

// ProvideBackupService creates and starts BackupService
func ProvideBackupService(
	settingRepo SettingRepository,
	cfg *config.Config,
	encryptor SecretEncryptor,
	storeFactory BackupObjectStoreFactory,
	dumper DBDumper,
) *BackupService {
	svc := NewBackupService(settingRepo, cfg, encryptor, storeFactory, dumper)
	svc.Start()
	return svc
}

// ProvideSettingService wires SettingService with group reader for default subscription validation.
func ProvideSettingService(settingRepo SettingRepository, groupRepo GroupRepository, cfg *config.Config) *SettingService {
	svc := NewSettingService(settingRepo, cfg)
	svc.SetDefaultSubscriptionGroupReader(groupRepo)
	return svc
}

// ProvidePaymentOrderExpiryService creates and starts PaymentOrderExpiryService.
func ProvidePaymentOrderExpiryService(
	orderRepo PaymentOrderRepository,
	auditLogRepo PaymentAuditLogRepository,
	registry *PaymentProviderRegistry,
	instanceRepo PaymentProviderInstanceRepository,
	encryptor SecretEncryptor,
) *PaymentOrderExpiryService {
	svc := NewPaymentOrderExpiryService(orderRepo, auditLogRepo, registry, instanceRepo, encryptor, 30*time.Second)
	svc.Start()
	return svc
}

// ProvideAPIKeyAuthCacheInvalidator provides a no-op cache invalidator for the base framework
func ProvideAPIKeyAuthCacheInvalidator() APIKeyAuthCacheInvalidator {
	return &noopAuthCacheInvalidator{}
}

type noopAuthCacheInvalidator struct{}

func (n *noopAuthCacheInvalidator) InvalidateAuthCacheByKey(_ context.Context, _ string)    {}
func (n *noopAuthCacheInvalidator) InvalidateAuthCacheByUserID(_ context.Context, _ int64)   {}
func (n *noopAuthCacheInvalidator) InvalidateAuthCacheByGroupID(_ context.Context, _ int64)  {}

// ProvideBillingCache provides a no-op billing cache for the base framework
func ProvideBillingCache() BillingCache {
	return &noopBillingCache{}
}

type noopBillingCache struct{}

func (n *noopBillingCache) GetUserBalance(_ context.Context, _ int64) (float64, error)          { return 0, nil }
func (n *noopBillingCache) SetUserBalance(_ context.Context, _ int64, _ float64) error           { return nil }
func (n *noopBillingCache) InvalidateUserBalance(_ context.Context, _ int64) error               { return nil }

// ProvideProxyService creates ProxyService and wires it with PaymentOrderService.
func ProvideProxyService(
	entClient *dbent.Client,
	nodeRepo ProxyNodeRepository,
	productRepo ProxyProductRepository,
	rentalRepo ProxyRentalRepository,
	credRepo ProxyCredentialRepository,
	trafficRepo ProxyTrafficLogRepository,
	paymentService *PaymentOrderService,
) *ProxyService {
	svc := NewProxyService(entClient, nodeRepo, productRepo, rentalRepo, credRepo, trafficRepo)
	paymentService.SetProxyService(svc)
	return svc
}

// ProvideProxyExpiryService creates and starts the proxy rental expiry background goroutine.
func ProvideProxyExpiryService(proxyService *ProxyService) *ProxyExpiryService {
	svc := NewProxyExpiryService(proxyService)
	svc.Start()
	return svc
}

// ProviderSet is the Wire provider set for all services
var ProviderSet = wire.NewSet(
	// Core services
	NewAuthService,
	NewUserService,
	ProvideAPIKeyAuthCacheInvalidator,
	ProvideBillingCache,
	NewBillingCacheService,
	NewRedeemService,
	NewPromoService,
	NewAnnouncementService,
	NewAdminService,
	NewEmailService,
	ProvideEmailQueueService,
	NewTurnstileService,
	NewSubscriptionService,
	wire.Bind(new(DefaultSubscriptionAssigner), new(*SubscriptionService)),
	ProvideUpdateService,
	ProvideSubscriptionExpiryService,
	NewUserAttributeService,
	NewTotpService,
	ProvideIdempotencyCoordinator,
	ProvideSystemOperationLockService,
	ProvideSettingService,
	ProvideBackupService,

	// Payment services
	NewPaymentConfigService,
	NewPaymentProviderRegistry,
	NewPaymentLoadBalancer,
	NewPaymentOrderService,
	ProvidePaymentOrderExpiryService,

	// Proxy services
	ProvideProxyService,
	ProvideProxyExpiryService,
)
