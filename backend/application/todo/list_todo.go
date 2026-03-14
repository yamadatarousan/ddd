package todo

import "errors"

import domain "github.com/user/ddd/backend/domain/todo"

var ErrCompletedQueryInvalid = errors.New("completedはtrueまたはfalseで指定してください")

// ListRepositoryは一覧取得ユースケースが必要とする最小インターフェース。
// 取得専用に絞ることで、読み取りユースケースの責務を明確に保つ。
type ListRepository interface {
	FindAll() ([]domain.Entity, error)
	FindByCompleted(isCompleted bool) ([]domain.Entity, error)
}

type ListTodoCommand struct {
	Completed *bool
}

type ListTodoUseCase struct {
	repository ListRepository
}

func NewListTodoUseCase(repository ListRepository) ListTodoUseCase {
	return ListTodoUseCase{repository: repository}
}

func (u ListTodoUseCase) Execute(command ListTodoCommand) ([]domain.Entity, error) {
	if command.Completed != nil {
		return u.repository.FindByCompleted(*command.Completed)
	}
	return u.repository.FindAll()
}
