package memory_test

import (
	"testing"
	"time"

	domain "github.com/user/ddd/backend/domain/notification"
	"github.com/user/ddd/backend/infrastructure/notification/memory"
)

func Test通知を保存して一覧取得できること(t *testing.T) {
	repository := memory.NewNotificationRepository()
	notification, _ := domain.NewNotification("notification-1", "Todo完了", time.Now())
	if err := repository.Save(notification); err != nil {
		t.Fatalf("保存に失敗: %v", err)
	}
	items, err := repository.FindAll()
	if err != nil {
		t.Fatalf("一覧取得に失敗: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("1件取得されるべき: got=%d", len(items))
	}
}
