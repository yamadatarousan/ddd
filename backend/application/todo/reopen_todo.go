package todo

import (
	"strings"

	domain "github.com/user/ddd/backend/domain/todo"
)

// ReopenRepositoryは未完了へ戻すユースケースが必要な永続化操作。
// 対象取得と保存のみ定義し、振る舞いの焦点を状態遷移に絞る。
type ReopenRepository interface {
	FindByID(id string) (domain.Entity, bool, error)
	Save(entity domain.Entity) error
}

type ReopenTodoCommand struct {
	ID string
}

type ReopenTodoUseCase struct {
	repository ReopenRepository
}

func NewReopenTodoUseCase(repository ReopenRepository) ReopenTodoUseCase {
	return ReopenTodoUseCase{repository: repository}
}

func (u ReopenTodoUseCase) Execute(command ReopenTodoCommand) (domain.Entity, error) {
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

	entity.Reopen()
	if err := u.repository.Save(entity); err != nil {
		return domain.Entity{}, err
	}

	return entity, nil
}
