package notification_test

import (
	"testing"
	"time"

	"github.com/user/ddd/backend/application/notification"
	domain "github.com/user/ddd/backend/domain/notification"
)

type メモリ通知リポジトリ struct {
	items []domain.Notification
}

type 固定ID生成器 struct{}

type 固定時計 struct{}

func (r *メモリ通知リポジトリ) Save(notification domain.Notification) error {
	r.items = append(r.items, notification)
	return nil
}

func (r *メモリ通知リポジトリ) FindAll() ([]domain.Notification, error) {
	return r.items, nil
}

func (g 固定ID生成器) NextID() string {
	return "notification-1"
}

func (c 固定時計) Now() time.Time {
	return time.Date(2026, 3, 1, 12, 0, 0, 0, time.UTC)
}

func TestTodo完了通知を記録できること(t *testing.T) {
	repository := &メモリ通知リポジトリ{}
	useCase := notification.NewRecordTodoCompletedUseCase(repository, 固定ID生成器{}, 固定時計{})

	notification, err := useCase.Execute(notification.RecordTodoCompletedCommand{TodoID: "todo-1", Title: "牛乳を買う"})
	if err != nil {
		t.Fatalf("通知記録は成功するべき: %v", err)
	}
	if notification.ID() != "notification-1" {
		t.Fatalf("通知IDが一致しない: got=%s", notification.ID())
	}
	if len(repository.items) != 1 {
		t.Fatalf("1件保存されるべき: got=%d", len(repository.items))
	}
}
