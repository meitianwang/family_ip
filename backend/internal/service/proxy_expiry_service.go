package service

import (
	"context"
	"log/slog"
	"time"
)

// ProxyExpiryService runs a background ticker to expire overdue proxy rentals.
type ProxyExpiryService struct {
	proxyService *ProxyService
	interval     time.Duration
	stop         chan struct{}
}

func NewProxyExpiryService(proxyService *ProxyService) *ProxyExpiryService {
	return &ProxyExpiryService{
		proxyService: proxyService,
		interval:     time.Hour,
		stop:         make(chan struct{}),
	}
}

func (s *ProxyExpiryService) Start() {
	go func() {
		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
				s.proxyService.ExpireOverdueRentals(ctx)
				cancel()
			case <-s.stop:
				slog.Info("ProxyExpiryService stopped")
				return
			}
		}
	}()
	slog.Info("ProxyExpiryService started", "interval", s.interval)
}

func (s *ProxyExpiryService) Stop() {
	close(s.stop)
}
