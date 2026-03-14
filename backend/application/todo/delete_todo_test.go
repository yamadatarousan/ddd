package todo_test

import (
	"testing"

	app "github.com/user/ddd/backend/application/todo"
	domain "github.com/user/ddd/backend/domain/todo"
)

type 削除用メモリリポジトリ struct {
	items map[string]domain.Entity
}

func new削除用メモリリポジトリ() *削除用メモリリポジトリ {
	return &削除用メモリリポジトリ{items: map[string]domain.Entity{}}
}

func (r *削除用メモリリポジトリ) Save(entity domain.Entity) error {
	r.items[entity.ID()] = entity
	return nil
}

func (r *削除用メモリリポジトリ) DeleteByID(id string) (bool, error) {
	if _, ok := r.items[id]; !ok {
		return false, nil
	}
	delete(r.items, id)
	return true, nil
}

func (r *削除用メモリリポジトリ) FindByID(id string) (domain.Entity, bool, error) {
	entity, ok := r.items[id]
	return entity, ok, nil
}

func Test存在するTodoを削除できること(t *testing.T) {
	repo := new削除用メモリリポジトリ()
	title, _ := domain.NewTitle("牛乳を買う")
	_ = repo.Save(domain.NewEntity("todo-1", title))

	useCase := app.NewDeleteTodoUseCase(repo)
	err := useCase.Execute(app.DeleteTodoCommand{ID: "todo-1"})
	if err != nil {
		t.Fatalf("削除は成功するべき: %v", err)
	}

	_, ok, _ := repo.FindByID("todo-1")
	if ok {
		t.Fatalf("削除後は存在しないべき")
	}
}

func Test存在しないTodoを削除しようとするとエラーになること(t *testing.T) {
	repo := new削除用メモリリポジトリ()
	useCase := app.NewDeleteTodoUseCase(repo)

	err := useCase.Execute(app.DeleteTodoCommand{ID: "todo-404"})
	if err == nil {
		t.Fatalf("存在しないTodoはエラーになるべき")
	}
}
