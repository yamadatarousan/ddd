package todo_test

import (
	"testing"

	app "github.com/user/ddd/backend/application/todo"
	domain "github.com/user/ddd/backend/domain/todo"
)

type 一覧用メモリリポジトリ struct {
	items []domain.Entity
}

func (r *一覧用メモリリポジトリ) Save(entity domain.Entity) error {
	r.items = append(r.items, entity)
	return nil
}

func (r *一覧用メモリリポジトリ) FindAll() ([]domain.Entity, error) {
	return r.items, nil
}

func (r *一覧用メモリリポジトリ) FindByCompleted(isCompleted bool) ([]domain.Entity, error) {
	filtered := make([]domain.Entity, 0, len(r.items))
	for _, item := range r.items {
		if item.IsCompleted() == isCompleted {
			filtered = append(filtered, item)
		}
	}
	return filtered, nil
}

func Test登録済みTodo一覧を取得できること(t *testing.T) {
	repo := &一覧用メモリリポジトリ{}
	title1, _ := domain.NewTitle("牛乳を買う")
	title2, _ := domain.NewTitle("散歩する")
	_ = repo.Save(domain.NewEntity("todo-1", title1))
	_ = repo.Save(domain.NewEntity("todo-2", title2))

	usecase := app.NewListTodoUseCase(repo)
	entities, err := usecase.Execute(app.ListTodoCommand{})
	if err != nil {
		t.Fatalf("一覧取得は成功するべき: %v", err)
	}
	if len(entities) != 2 {
		t.Fatalf("2件取得できるべき: got=%d", len(entities))
	}
	if entities[0].ID() != "todo-1" || entities[1].ID() != "todo-2" {
		t.Fatalf("登録順で取得できるべき")
	}
}

func Test完了状態を指定して一覧取得できること(t *testing.T) {
	repo := &一覧用メモリリポジトリ{}
	title1, _ := domain.NewTitle("牛乳を買う")
	title2, _ := domain.NewTitle("散歩する")

	completed := domain.NewEntity("todo-1", title1)
	completed.Complete()
	_ = repo.Save(completed)
	_ = repo.Save(domain.NewEntity("todo-2", title2))

	usecase := app.NewListTodoUseCase(repo)
	trueValue := true
	entities, err := usecase.Execute(app.ListTodoCommand{Completed: &trueValue})
	if err != nil {
		t.Fatalf("一覧取得は成功するべき: %v", err)
	}
	if len(entities) != 1 {
		t.Fatalf("1件取得できるべき: got=%d", len(entities))
	}
	if entities[0].ID() != "todo-1" {
		t.Fatalf("完了Todoのみ取得されるべき: got=%s", entities[0].ID())
	}
}
