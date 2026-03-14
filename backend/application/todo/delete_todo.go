package todo

import "strings"

// DeleteRepositoryは削除ユースケースが必要とする永続化操作。
// 削除対象の有無をboolで受け取り、存在しない場合を明示的に扱えるようにする。
type DeleteRepository interface {
	DeleteByID(id string) (bool, error)
}

type DeleteTodoCommand struct {
	ID string
}

type DeleteTodoUseCase struct {
	repository DeleteRepository
}

func NewDeleteTodoUseCase(repository DeleteRepository) DeleteTodoUseCase {
	return DeleteTodoUseCase{repository: repository}
}

func (u DeleteTodoUseCase) Execute(command DeleteTodoCommand) error {
	id := strings.TrimSpace(command.ID)
	if id == "" {
		return ErrTodoIDRequired
	}

	deleted, err := u.repository.DeleteByID(id)
	if err != nil {
		return err
	}
	if !deleted {
		return ErrTodoNotFound
	}
	return nil
}
