package main

import (
	"fmt"
	"log"
	"os"
	"sync/atomic"

	notification "github.com/user/ddd/backend/application/notification"
	app "github.com/user/ddd/backend/application/todo"
	notificationmemory "github.com/user/ddd/backend/infrastructure/notification/memory"
	notificationsystem "github.com/user/ddd/backend/infrastructure/notification/system"
	"github.com/user/ddd/backend/infrastructure/todo/memory"
	"github.com/user/ddd/backend/integration/todo_notification"
	httpapi "github.com/user/ddd/backend/interfaces/http"
)

func main() {
	todoRepository := memory.NewTodoRepository()
	notificationRepository := notificationmemory.NewNotificationRepository()

	var idSequence uint64
	generateID := func() string {
		id := atomic.AddUint64(&idSequence, 1)
		return fmt.Sprintf("todo-%d", id)
	}

	recordNotificationUseCase := notification.NewRecordTodoCompletedUseCase(
		notificationRepository,
		notificationsystem.NewSequenceIDGenerator(),
		notificationsystem.NewRealtimeClock(),
	)
	listNotificationUseCase := notification.NewListNotificationUseCase(notificationRepository)
	notifier := todo_notification.NewNotifier(recordNotificationUseCase)

	createUseCase := app.NewCreateTodoUseCase(todoRepository, generateID)
	completeUseCase := app.NewCompleteTodoUseCaseWithNotifier(todoRepository, notifier)
	listUseCase := app.NewListTodoUseCase(todoRepository)
	updateTitleUseCase := app.NewUpdateTodoTitleUseCase(todoRepository)
	deleteUseCase := app.NewDeleteTodoUseCase(todoRepository)
	reopenUseCase := app.NewReopenTodoUseCase(todoRepository)

	router := httpapi.NewRouter(
		createUseCase,
		completeUseCase,
		listUseCase,
		updateTitleUseCase,
		deleteUseCase,
		reopenUseCase,
		listNotificationUseCase,
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 学習用途の最小サーバーなので、起動失敗時は即終了して原因を明示する。
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("サーバー起動に失敗しました: %v", err)
	}
}
