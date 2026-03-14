package todo_test

import (
	"testing"

	app "github.com/user/ddd/backend/application/todo"
	domain "github.com/user/ddd/backend/domain/todo"
)

type 完了用メモリリポジトリ struct {
	items map[string]domain.Entity
}

func new完了用メモリリポジトリ() *完了用メモリリポジトリ {
	return &完了用メモリリポジトリ{items: map[string]domain.Entity{}}
}

func (r *完了用メモリリポジトリ) Save(entity domain.Entity) error {
	r.items[entity.ID()] = entity
	return nil
}

func (r *完了用メモリリポジトリ) FindByID(id string) (domain.Entity, bool, error) {
	entity, ok := r.items[id]
	return entity, ok, nil
}

func Test未完了Todoを完了にできること(t *testing.T) {
	repo := new完了用メモリリポジトリ()
	title, _ := domain.NewTitle("牛乳を買う")
	_ = repo.Save(domain.NewEntity("todo-1", title))

	usecase := app.NewCompleteTodoUseCase(repo)
	entity, err := usecase.Execute(app.CompleteTodoCommand{ID: "todo-1"})
	if err != nil {
		t.Fatalf("完了処理は成功するべき: %v", err)
	}
	if !entity.IsCompleted() {
		t.Fatalf("完了状態であるべき")
	}

	saved, ok, _ := repo.FindByID("todo-1")
	if !ok {
		t.Fatalf("保存済みTodoが見つからない")
	}
	if !saved.IsCompleted() {
		t.Fatalf("保存後も完了状態であるべき")
	}
}

func Test存在しないTodoを完了しようとするとエラーになること(t *testing.T) {
	repo := new完了用メモリリポジトリ()
	usecase := app.NewCompleteTodoUseCase(repo)

	_, err := usecase.Execute(app.CompleteTodoCommand{ID: "todo-404"})
	if err == nil {
		t.Fatalf("存在しないTodoはエラーになるべき")
	}
}
