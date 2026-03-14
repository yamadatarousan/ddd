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

type CompleteTodoCommand struct {
	ID string
}

type CompleteTodoUseCase struct {
	repository CompleteRepository
}

func NewCompleteTodoUseCase(repository CompleteRepository) CompleteTodoUseCase {
	return CompleteTodoUseCase{repository: repository}
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

	return entity, nil
}
