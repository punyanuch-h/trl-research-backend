package cron

import (
	"log"
	"trl-research-backend/internal/repository"
)

type TokenCleanupCron struct {
	RefreshTokenRepo repository.RefreshTokenRepository
}

func NewTokenCleanupCron(repo repository.RefreshTokenRepository) *TokenCleanupCron {
	return &TokenCleanupCron{
		RefreshTokenRepo: repo,
	}
}

func (c *TokenCleanupCron) Start() {
	// We use the same central robfig/cron/v3 approach or a simple ticker
	// For consistency with ReminderCron, let's assume we can add it to a shared cron.
	// But since ReminderCron creates its own instance, we'll do the same here for modularity.

	cr := NewCronWithTimezone()

	// Scheduling: Run every day at midnight
	// 0 0 * * *
	_, err := cr.AddFunc("0 0 * * *", func() {
		log.Println("🧹 Running refresh token cleanup job...")
		if err := c.RefreshTokenRepo.DeleteExpired(); err != nil {
			log.Printf("❌ Error cleaning up expired tokens: %v", err)
		} else {
			log.Println("✅ Expired tokens cleaned up successfully")
		}
	})

	if err != nil {
		log.Fatalf("❌ Error adding cleanup cron job: %v", err)
	}

	cr.Start()
	log.Println("🚀 Token cleanup cron job started - Running daily at midnight")
}
