package todo_notification_test

import (
	"testing"
	"time"

	notification "github.com/user/ddd/backend/application/notification"
	domain "github.com/user/ddd/backend/domain/notification"
	"github.com/user/ddd/backend/integration/todo_notification"
)

type テスト通知リポジトリ struct {
	items []domain.Notification
}

type テストID生成器 struct{}

type テスト時計 struct{}

func (r *テスト通知リポジトリ) Save(notification domain.Notification) error {
	r.items = append(r.items, notification)
	return nil
}

func (r *テスト通知リポジトリ) FindAll() ([]domain.Notification, error) {
	return r.items, nil
}

func (g テストID生成器) NextID() string {
	return "notification-1"
}

func (c テスト時計) Now() time.Time {
	return time.Date(2026, 3, 1, 12, 0, 0, 0, time.UTC)
}

func Test完了通知を通知コンテキストに橋渡しできること(t *testing.T) {
	repository := &テスト通知リポジトリ{}
	recordUseCase := notification.NewRecordTodoCompletedUseCase(repository, テストID生成器{}, テスト時計{})
	notifier := todo_notification.NewNotifier(recordUseCase)

	if err := notifier.NotifyTodoCompleted("todo-1", "牛乳を買う"); err != nil {
		t.Fatalf("通知橋渡しは成功するべき: %v", err)
	}
	if len(repository.items) != 1 {
		t.Fatalf("通知は1件記録されるべき: got=%d", len(repository.items))
	}
}
