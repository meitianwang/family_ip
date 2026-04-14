//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/meitianwang/fast-frame/ent"
	"github.com/meitianwang/fast-frame/internal/config"
	"github.com/meitianwang/fast-frame/internal/handler"
	"github.com/meitianwang/fast-frame/internal/repository"
	"github.com/meitianwang/fast-frame/internal/server"
	"github.com/meitianwang/fast-frame/internal/server/middleware"
	"github.com/meitianwang/fast-frame/internal/service"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

type Application struct {
	Server  *http.Server
	Cleanup func()
}

func initializeApplication(buildInfo handler.BuildInfo) (*Application, error) {
	wire.Build(
		// Infrastructure layer ProviderSets
		config.ProviderSet,

		// Business layer ProviderSets
		repository.ProviderSet,
		service.ProviderSet,
		middleware.ProviderSet,
		handler.ProviderSet,

		// Server layer ProviderSet
		server.ProviderSet,

		// BuildInfo provider
		provideServiceBuildInfo,

		// Cleanup function provider
		provideCleanup,

		// Application struct
		wire.Struct(new(Application), "Server", "Cleanup"),
	)
	return nil, nil
}

func provideServiceBuildInfo(buildInfo handler.BuildInfo) service.BuildInfo {
	return service.BuildInfo{
		Version:   buildInfo.Version,
		BuildType: buildInfo.BuildType,
	}
}

func provideCleanup(
	entClient *ent.Client,
	rdb *redis.Client,
	subscriptionExpiry *service.SubscriptionExpiryService,
	emailQueue *service.EmailQueueService,
	subscriptionService *service.SubscriptionService,
	backupSvc *service.BackupService,
	paymentOrderExpiry *service.PaymentOrderExpiryService,
	proxyExpiry *service.ProxyExpiryService,
) func() {
	return func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Stop application services
		if subscriptionExpiry != nil {
			subscriptionExpiry.Stop()
			log.Println("[Cleanup] SubscriptionExpiryService stopped")
		}
		if subscriptionService != nil {
			subscriptionService.Stop()
			log.Println("[Cleanup] SubscriptionService stopped")
		}
		if emailQueue != nil {
			emailQueue.Stop()
			log.Println("[Cleanup] EmailQueueService stopped")
		}
		if backupSvc != nil {
			backupSvc.Stop()
			log.Println("[Cleanup] BackupService stopped")
		}
		if paymentOrderExpiry != nil {
			paymentOrderExpiry.Stop()
			log.Println("[Cleanup] PaymentOrderExpiryService stopped")
		}
		if proxyExpiry != nil {
			proxyExpiry.Stop()
			log.Println("[Cleanup] ProxyExpiryService stopped")
		}

		// Close infrastructure
		if rdb != nil {
			if err := rdb.Close(); err != nil {
				log.Printf("[Cleanup] Redis close failed: %v", err)
			} else {
				log.Println("[Cleanup] Redis closed")
			}
		}
		if entClient != nil {
			if err := entClient.Close(); err != nil {
				log.Printf("[Cleanup] Ent close failed: %v", err)
			} else {
				log.Println("[Cleanup] Ent closed")
			}
		}

		select {
		case <-ctx.Done():
			log.Println("[Cleanup] Warning: cleanup timed out after 10 seconds")
		default:
			log.Println("[Cleanup] All cleanup steps completed")
		}
	}
}
