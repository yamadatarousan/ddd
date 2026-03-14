package notification

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrNotificationIDRequired      = errors.New("通知IDは必須です")
	ErrNotificationMessageRequired = errors.New("通知メッセージは必須です")
	ErrNotificationCreatedAtEmpty  = errors.New("通知作成日時は必須です")
)

// Notificationは通知コンテキストのエンティティ。
// Todo管理とは別のモデルとして、通知表示に必要な情報のみを持つ。
type Notification struct {
	id        string
	message   string
	createdAt time.Time
	isRead    bool
}

func NewNotification(id string, message string, createdAt time.Time) (Notification, error) {
	trimmedID := strings.TrimSpace(id)
	if trimmedID == "" {
		return Notification{}, ErrNotificationIDRequired
	}
	trimmedMessage := strings.TrimSpace(message)
	if trimmedMessage == "" {
		return Notification{}, ErrNotificationMessageRequired
	}
	if createdAt.IsZero() {
		return Notification{}, ErrNotificationCreatedAtEmpty
	}

	return Notification{
		id:        trimmedID,
		message:   trimmedMessage,
		createdAt: createdAt,
		isRead:    false,
	}, nil
}

func (n Notification) ID() string {
	return n.id
}

func (n Notification) Message() string {
	return n.message
}

func (n Notification) CreatedAt() time.Time {
	return n.createdAt
}

func (n Notification) IsRead() bool {
	return n.isRead
}

func (n *Notification) MarkRead() {
	n.isRead = true
}
