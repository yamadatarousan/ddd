package todo

import domain "github.com/user/ddd/backend/domain/todo"

// ListRepositoryは一覧取得ユースケースが必要とする最小インターフェース。
// 取得専用に絞ることで、読み取りユースケースの責務を明確に保つ。
type ListRepository interface {
	FindAll() ([]domain.Entity, error)
}

type ListTodoUseCase struct {
	repository ListRepository
}

func NewListTodoUseCase(repository ListRepository) ListTodoUseCase {
	return ListTodoUseCase{repository: repository}
}

func (u ListTodoUseCase) Execute() ([]domain.Entity, error) {
	return u.repository.FindAll()
}
