package todo

import (
	"errors"
	"strings"

	domain "github.com/user/ddd/backend/domain/todo"
)

var (
	ErrTodoIDRequired = errors.New("Todo IDは必須です")
	ErrTodoNotFound   = errors.New("Todoが見つかりません")
)

// CompleteRepositoryは完了ユースケースが必要とする永続化操作を定義する。
// 読み出しと保存を最小限に限定し、ユースケースの関心を明確にする。
type CompleteRepository interface {
	FindByID(id string) (domain.Entity, bool, error)
	Save(entity domain.Entity) error
}

// TodoCompletionNotifierは、Todo管理コンテキスト外へ完了通知を伝えるポート。
// 外部連携の詳細は持ち込まず、完了事実のみを通知する。
type TodoCompletionNotifier interface {
	NotifyTodoCompleted(todoID string, title string) error
}

type CompleteTodoCommand struct {
	ID string
}

type CompleteTodoUseCase struct {
	repository CompleteRepository
	notifier   TodoCompletionNotifier
}

type noopTodoCompletionNotifier struct{}

func (n noopTodoCompletionNotifier) NotifyTodoCompleted(_ string, _ string) error {
	return nil
}

func NewCompleteTodoUseCase(repository CompleteRepository) CompleteTodoUseCase {
	return CompleteTodoUseCase{
		repository: repository,
		notifier:   noopTodoCompletionNotifier{},
	}
}

func NewCompleteTodoUseCaseWithNotifier(
	repository CompleteRepository,
	notifier TodoCompletionNotifier,
) CompleteTodoUseCase {
	if notifier == nil {
		notifier = noopTodoCompletionNotifier{}
	}
	return CompleteTodoUseCase{
		repository: repository,
		notifier:   notifier,
	}
}

func (u CompleteTodoUseCase) Execute(command CompleteTodoCommand) (domain.Entity, error) {
	id := strings.TrimSpace(command.ID)
	if id == "" {
		return domain.Entity{}, ErrTodoIDRequired
	}

	entity, ok, err := u.repository.FindByID(id)
	if err != nil {
		return domain.Entity{}, err
	}
	if !ok {
		return domain.Entity{}, ErrTodoNotFound
	}

	entity.Complete()
	if err := u.repository.Save(entity); err != nil {
		return domain.Entity{}, err
	}
	if err := u.notifier.NotifyTodoCompleted(entity.ID(), entity.Title().Value()); err != nil {
		return domain.Entity{}, err
	}

	return entity, nil
}
