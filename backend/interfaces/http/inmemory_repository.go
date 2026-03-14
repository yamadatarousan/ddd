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

	// 同一IDが既にある場合は置き換え、なければ追加する。
	// 完了ユースケースで状態更新を扱うため、追記のみではなく上書きを許可する。
	for index := range r.items {
		if r.items[index].ID() == entity.ID() {
			r.items[index] = entity
			return nil
		}
	}

	r.items = append(r.items, entity)
	return nil
}

func (r *InMemoryTodoRepository) FindByID(id string) (domain.Entity, bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, item := range r.items {
		if item.ID() == id {
			return item, true, nil
		}
	}

	return domain.Entity{}, false, nil
}

func (r *InMemoryTodoRepository) FindAll() ([]domain.Entity, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 内部配列の参照をそのまま返すと呼び出し側が更新できるため、コピーを返す。
	copied := make([]domain.Entity, len(r.items))
	copy(copied, r.items)
	return copied, nil
}
