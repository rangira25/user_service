package service

import (
	"log"
	"time"
)

func logAudit(action, userID string) {
	log.Printf("[AUDIT] %s on user %s at %s", action, userID, time.Now().Format(time.RFC3339))
}
