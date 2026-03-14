package main

import (
	"fmt"
	"log"
	"os"
	"sync/atomic"

	app "github.com/user/ddd/backend/application/todo"
	httpapi "github.com/user/ddd/backend/interfaces/http"
)

func main() {
	repository := httpapi.NewInMemoryTodoRepository()

	var idSequence uint64
	generateID := func() string {
		id := atomic.AddUint64(&idSequence, 1)
		return fmt.Sprintf("todo-%d", id)
	}

	createUseCase := app.NewCreateTodoUseCase(repository, generateID)
	completeUseCase := app.NewCompleteTodoUseCase(repository)
	listUseCase := app.NewListTodoUseCase(repository)
	updateTitleUseCase := app.NewUpdateTodoTitleUseCase(repository)
	deleteUseCase := app.NewDeleteTodoUseCase(repository)

	router := httpapi.NewRouter(
		createUseCase,
		completeUseCase,
		listUseCase,
		updateTitleUseCase,
		deleteUseCase,
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
