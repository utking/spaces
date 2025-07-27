package ports

import (
	"context"

	"gogs.utking.net/utking/spaces/internal/application/domain"
)

// NotificationService is an interface that defines the methods for sending notifications.
type NotificationService interface {
	Send(ctx context.Context, message *domain.Notification) error
}
