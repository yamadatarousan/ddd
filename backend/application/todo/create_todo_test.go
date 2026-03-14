package todo_test

import (
	"testing"

	app "github.com/user/ddd/backend/application/todo"
	domain "github.com/user/ddd/backend/domain/todo"
)

type メモリリポジトリ struct {
	saved []domain.Entity
}

func (r *メモリリポジトリ) Save(entity domain.Entity) error {
	r.saved = append(r.saved, entity)
	return nil
}

func Testタイトルを指定してTodoを作成して保存できること(t *testing.T) {
	repo := &メモリリポジトリ{}
	usecase := app.NewCreateTodoUseCase(repo, func() string { return "todo-1" })

	entity, err := usecase.Execute(app.CreateTodoCommand{Title: "牛乳を買う"})
	if err != nil {
		t.Fatalf("作成は成功するべき: %v", err)
	}

	if entity.ID() != "todo-1" {
		t.Fatalf("IDが一致しない: got=%s", entity.ID())
	}
	if len(repo.saved) != 1 {
		t.Fatalf("1件保存されるべき: got=%d", len(repo.saved))
	}
	if repo.saved[0].Title().Value() != "牛乳を買う" {
		t.Fatalf("保存されたタイトルが一致しない: got=%s", repo.saved[0].Title().Value())
	}
}

func Test不正なタイトルの場合はエラーになること(t *testing.T) {
	repo := &メモリリポジトリ{}
	usecase := app.NewCreateTodoUseCase(repo, func() string { return "todo-1" })

	_, err := usecase.Execute(app.CreateTodoCommand{Title: ""})
	if err == nil {
		t.Fatalf("不正タイトルはエラーになるべき")
	}
	if len(repo.saved) != 0 {
		t.Fatalf("エラー時は保存されないべき: got=%d", len(repo.saved))
	}
}
