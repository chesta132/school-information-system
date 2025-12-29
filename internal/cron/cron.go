package cron

import (
	"context"
	"log"
	"school-information-system/config"
	"school-information-system/internal/repos"
	"time"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	_ "time/tzdata"
)

func New(db *gorm.DB) *cron.Cron {
	userRepo := repos.NewUser(db)
	revokedRepo := repos.NewRevoked(db)
	ctx := context.Background()

	loc, err := time.LoadLocation(config.APP_TIMEZONE)
	if err != nil {
		log.Printf("[CRON] Failed to load timezone, using UTC: %v", err)
		loc = time.UTC
	}
	scheduler := cron.New(cron.WithLocation(loc))

	// Delete soft-deleted users that have been deleted for more than configured duration
	scheduler.AddFunc(config.CRON_USER_INTERVAL, func() {
		rows, err := userRepo.DeleteAllUser(ctx, "deleted_at IS NOT NULL AND deleted_at <= ?", time.Now().Add(-config.CRON_USER_DELETE_AFTER))

		if err != nil {
			log.Printf("[CLEANUP] Failed to delete soft-deleted users: %v", err)
			return
		}

		if rows > 0 {
			// TODO: send notification to admins about deleted users
			log.Printf("[CLEANUP] Deleted %d soft-deleted users at %s", rows, time.Now().Format(time.RFC3339))
		}
	})

	// Delete revoked tokens that have expired
	scheduler.AddFunc(config.CRON_REVOKED_INTERVAL, func() {
		rows, err := revokedRepo.Delete(ctx, "revoked_until IS NOT NULL AND revoked_until <= ?", time.Now())

		if err != nil {
			log.Printf("[CLEANUP] Failed to delete revoked tokens: %v", err)
			return
		}

		if rows > 0 {
			log.Printf("[CLEANUP] Deleted %d revoked tokens at %s", rows, time.Now().Format(time.RFC3339))
		}
	})

	return scheduler
}
