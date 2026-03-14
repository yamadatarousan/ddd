package memory

import (
	"sync"

	domain "github.com/user/ddd/backend/domain/todo"
)

// TodoRepositoryは検証・学習用のインメモリ実装。
// 複数ユースケースから同時に参照される前提で排他制御を行う。
type TodoRepository struct {
	mu    sync.Mutex
	items []domain.Entity
}

func NewTodoRepository() *TodoRepository {
	return &TodoRepository{items: make([]domain.Entity, 0)}
}

func (r *TodoRepository) Save(entity domain.Entity) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for index := range r.items {
		if r.items[index].ID() == entity.ID() {
			r.items[index] = entity
			return nil
		}
	}

	r.items = append(r.items, entity)
	return nil
}

func (r *TodoRepository) FindByID(id string) (domain.Entity, bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, item := range r.items {
		if item.ID() == id {
			return item, true, nil
		}
	}

	return domain.Entity{}, false, nil
}

func (r *TodoRepository) FindAll() ([]domain.Entity, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	copied := make([]domain.Entity, len(r.items))
	copy(copied, r.items)
	return copied, nil
}

func (r *TodoRepository) DeleteByID(id string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for index := range r.items {
		if r.items[index].ID() == id {
			r.items = append(r.items[:index], r.items[index+1:]...)
			return true, nil
		}
	}

	return false, nil
}
