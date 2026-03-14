package todo_test

import (
	"testing"

	app "github.com/user/ddd/backend/application/todo"
	domain "github.com/user/ddd/backend/domain/todo"
)

type 再開用メモリリポジトリ struct {
	items map[string]domain.Entity
}

func new再開用メモリリポジトリ() *再開用メモリリポジトリ {
	return &再開用メモリリポジトリ{items: map[string]domain.Entity{}}
}

func (r *再開用メモリリポジトリ) Save(entity domain.Entity) error {
	r.items[entity.ID()] = entity
	return nil
}

func (r *再開用メモリリポジトリ) FindByID(id string) (domain.Entity, bool, error) {
	entity, ok := r.items[id]
	return entity, ok, nil
}

func Test完了済みTodoを未完了に戻せること(t *testing.T) {
	repo := new再開用メモリリポジトリ()
	title, _ := domain.NewTitle("牛乳を買う")
	entity := domain.NewEntity("todo-1", title)
	entity.Complete()
	_ = repo.Save(entity)

	useCase := app.NewReopenTodoUseCase(repo)
	updated, err := useCase.Execute(app.ReopenTodoCommand{ID: "todo-1"})
	if err != nil {
		t.Fatalf("未完了への戻しは成功するべき: %v", err)
	}
	if updated.IsCompleted() {
		t.Fatalf("未完了状態であるべき")
	}
}

func Test存在しないTodoの未完了戻しはエラーになること(t *testing.T) {
	repo := new再開用メモリリポジトリ()
	useCase := app.NewReopenTodoUseCase(repo)

	_, err := useCase.Execute(app.ReopenTodoCommand{ID: "todo-404"})
	if err == nil {
		t.Fatalf("存在しないTodoはエラーになるべき")
	}
}
