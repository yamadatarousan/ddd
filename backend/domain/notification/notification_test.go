package notification_test

import (
	"testing"
	"time"

	notification "github.com/user/ddd/backend/domain/notification"
)

func Test有効な通知を作成できること(t *testing.T) {
	entity, err := notification.NewNotification("notification-1", "Todoを完了しました", time.Now())
	if err != nil {
		t.Fatalf("通知は作成できるべき: %v", err)
	}
	if entity.ID() != "notification-1" {
		t.Fatalf("IDが一致しない: got=%s", entity.ID())
	}
	if entity.IsRead() {
		t.Fatalf("作成時は未読であるべき")
	}
}

func Test空メッセージでは通知を作成できないこと(t *testing.T) {
	_, err := notification.NewNotification("notification-1", "", time.Now())
	if err == nil {
		t.Fatalf("空メッセージはエラーになるべき")
	}
}
