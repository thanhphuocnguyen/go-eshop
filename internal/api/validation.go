package api

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

// ValidateServerDependencies performs comprehensive validation of server dependencies
func (s *Server) ValidateServerDependencies(ctx context.Context) error {
	log.Info().Msg("Starting server dependency validation...")

	// Validate server instance
	if s == nil {
		return fmt.Errorf("server instance is nil")
	}

	// Validate repository
	if s.repo == nil {
		return fmt.Errorf("database repository is nil")
	}

	// Test database connection with a simple query
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Try a simple count query to test database connectivity
	_, err := s.repo.CountUsers(ctx)
	if err != nil {
		return fmt.Errorf("database connectivity test failed: %w", err)
	}
	log.Info().Msg("Database connection validated")

	// Validate other dependencies
	if s.tokenAuth == nil {
		return fmt.Errorf("token auth is nil")
	}
	log.Info().Msg("Token auth validated")

	if s.validator == nil {
		return fmt.Errorf("validator is nil")
	}
	log.Info().Msg("Validator validated")

	if s.uploadService == nil {
		return fmt.Errorf("upload service is nil")
	}
	log.Info().Msg("Upload service validated")

	if s.paymentSrv == nil {
		return fmt.Errorf("payment service is nil")
	}
	log.Info().Msg("Payment service validated")

	if s.cacheSrv == nil {
		return fmt.Errorf("cache service is nil")
	}
	log.Info().Msg("Cache service validated")

	if s.taskDistributor == nil {
		return fmt.Errorf("task distributor is nil")
	}
	log.Info().Msg("Task distributor validated")

	if s.discountProcessor == nil {
		return fmt.Errorf("discount processor is nil")
	}
	log.Info().Msg("Discount processor validated")

	if s.router == nil {
		return fmt.Errorf("router is nil")
	}
	log.Info().Msg("Router validated")

	log.Info().Msg("All server dependencies validated successfully")
	return nil
}
