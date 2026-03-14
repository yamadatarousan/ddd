package todo_notification

import notification "github.com/user/ddd/backend/application/notification"

// NotifierはTodo管理コンテキストから通知コンテキストへ橋渡しするアダプター。
type Notifier struct {
	recordUseCase notification.RecordTodoCompletedUseCase
}

func NewNotifier(recordUseCase notification.RecordTodoCompletedUseCase) Notifier {
	return Notifier{recordUseCase: recordUseCase}
}

func (n Notifier) NotifyTodoCompleted(todoID string, title string) error {
	_, err := n.recordUseCase.Execute(notification.RecordTodoCompletedCommand{TodoID: todoID, Title: title})
	return err
}
