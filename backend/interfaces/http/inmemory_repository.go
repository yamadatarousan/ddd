package http

import (
	"sync"

	domain "github.com/user/ddd/backend/domain/todo"
)

// InMemoryTodoRepositoryは検証用の最小実装。
// 永続化先を導入する前に、ユースケースとAPIの接続確認を可能にする。
type InMemoryTodoRepository struct {
	mu    sync.Mutex
	items []domain.Entity
}

func NewInMemoryTodoRepository() *InMemoryTodoRepository {
	return &InMemoryTodoRepository{items: make([]domain.Entity, 0)}
}

func (r *InMemoryTodoRepository) Save(entity domain.Entity) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.items = append(r.items, entity)
	return nil
}
