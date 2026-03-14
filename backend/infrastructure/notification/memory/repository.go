package memory

import (
	"sync"

	domain "github.com/user/ddd/backend/domain/notification"
)

// NotificationRepositoryは通知コンテキスト用のインメモリ実装。
type NotificationRepository struct {
	mu    sync.Mutex
	items []domain.Notification
}

func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{items: make([]domain.Notification, 0)}
}

func (r *NotificationRepository) Save(notification domain.Notification) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.items = append(r.items, notification)
	return nil
}

func (r *NotificationRepository) FindAll() ([]domain.Notification, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	copied := make([]domain.Notification, len(r.items))
	copy(copied, r.items)
	return copied, nil
}
