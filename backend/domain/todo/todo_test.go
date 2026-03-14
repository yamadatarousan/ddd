package todo_test

import (
	"testing"

	"github.com/user/ddd/backend/domain/todo"
)

const 超過タイトル = "123456789012345678901234567890123456789012345678901234567890123456789012345678901"

func Testタイトルが空文字の場合は作成できないこと(t *testing.T) {
	_, err := todo.NewTitle("")
	if err == nil {
		t.Fatalf("空文字はエラーになるべき")
	}
}

func Test80文字を超えるタイトルは作成できないこと(t *testing.T) {
	_, err := todo.NewTitle(超過タイトル)
	if err == nil {
		t.Fatalf("80文字を超えるタイトルはエラーになるべき")
	}
}

func Test有効なタイトルからTodoを作成できること(t *testing.T) {
	title, err := todo.NewTitle("牛乳を買う")
	if err != nil {
		t.Fatalf("有効なタイトルは作成できるべき: %v", err)
	}

	entity := todo.NewEntity("todo-1", title)

	if entity.ID() != "todo-1" {
		t.Fatalf("IDが一致しない: got=%s", entity.ID())
	}
	if entity.Title().Value() != "牛乳を買う" {
		t.Fatalf("Titleが一致しない: got=%s", entity.Title().Value())
	}
	if entity.IsCompleted() {
		t.Fatalf("新規作成時は未完了であるべき")
	}
}

func Test完了に変更できること(t *testing.T) {
	title, _ := todo.NewTitle("牛乳を買う")
	entity := todo.NewEntity("todo-1", title)

	entity.Complete()

	if !entity.IsCompleted() {
		t.Fatalf("Complete後は完了状態であるべき")
	}
}

func Test完了済みを再度完了にしても状態が壊れないこと(t *testing.T) {
	title, _ := todo.NewTitle("牛乳を買う")
	entity := todo.NewEntity("todo-1", title)

	entity.Complete()
	entity.Complete()

	if !entity.IsCompleted() {
		t.Fatalf("再実行しても完了状態を維持するべき")
	}
}
