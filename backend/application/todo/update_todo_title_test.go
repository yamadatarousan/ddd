package todo_test

import (
	"testing"

	app "github.com/user/ddd/backend/application/todo"
	domain "github.com/user/ddd/backend/domain/todo"
)

type タイトル変更用メモリリポジトリ struct {
	items map[string]domain.Entity
}

func newタイトル変更用メモリリポジトリ() *タイトル変更用メモリリポジトリ {
	return &タイトル変更用メモリリポジトリ{items: map[string]domain.Entity{}}
}

func (r *タイトル変更用メモリリポジトリ) FindByID(id string) (domain.Entity, bool, error) {
	entity, ok := r.items[id]
	return entity, ok, nil
}

func (r *タイトル変更用メモリリポジトリ) Save(entity domain.Entity) error {
	r.items[entity.ID()] = entity
	return nil
}

func TestTodoのタイトルを変更できること(t *testing.T) {
	repo := newタイトル変更用メモリリポジトリ()
	title, _ := domain.NewTitle("牛乳を買う")
	_ = repo.Save(domain.NewEntity("todo-1", title))

	useCase := app.NewUpdateTodoTitleUseCase(repo)
	entity, err := useCase.Execute(app.UpdateTodoTitleCommand{ID: "todo-1", Title: "パンを買う"})
	if err != nil {
		t.Fatalf("タイトル変更は成功するべき: %v", err)
	}
	if entity.Title().Value() != "パンを買う" {
		t.Fatalf("タイトルが更新されていない: got=%s", entity.Title().Value())
	}
}

func Test不正なタイトルへ変更するとエラーになること(t *testing.T) {
	repo := newタイトル変更用メモリリポジトリ()
	title, _ := domain.NewTitle("牛乳を買う")
	_ = repo.Save(domain.NewEntity("todo-1", title))

	useCase := app.NewUpdateTodoTitleUseCase(repo)
	_, err := useCase.Execute(app.UpdateTodoTitleCommand{ID: "todo-1", Title: ""})
	if err == nil {
		t.Fatalf("不正タイトルはエラーになるべき")
	}
}

func Test存在しないTodoのタイトル変更はエラーになること(t *testing.T) {
	repo := newタイトル変更用メモリリポジトリ()
	useCase := app.NewUpdateTodoTitleUseCase(repo)

	_, err := useCase.Execute(app.UpdateTodoTitleCommand{ID: "todo-404", Title: "パンを買う"})
	if err == nil {
		t.Fatalf("存在しないTodoはエラーになるべき")
	}
}
