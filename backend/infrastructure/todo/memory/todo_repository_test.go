package memory_test

import (
	"testing"

	domain "github.com/user/ddd/backend/domain/todo"
	"github.com/user/ddd/backend/infrastructure/todo/memory"
)

func Test保存したTodoをIDで取得できること(t *testing.T) {
	repository := memory.NewTodoRepository()
	title, _ := domain.NewTitle("牛乳を買う")
	saved := domain.NewEntity("todo-1", title)
	if err := repository.Save(saved); err != nil {
		t.Fatalf("保存に失敗: %v", err)
	}

	entity, ok, err := repository.FindByID("todo-1")
	if err != nil {
		t.Fatalf("取得に失敗: %v", err)
	}
	if !ok {
		t.Fatalf("保存済みTodoが見つからない")
	}
	if entity.Title().Value() != "牛乳を買う" {
		t.Fatalf("タイトルが一致しない: got=%s", entity.Title().Value())
	}
}

func TestFindAllは内部状態のコピーを返すこと(t *testing.T) {
	repository := memory.NewTodoRepository()
	title, _ := domain.NewTitle("牛乳を買う")
	_ = repository.Save(domain.NewEntity("todo-1", title))

	entities, err := repository.FindAll()
	if err != nil {
		t.Fatalf("一覧取得に失敗: %v", err)
	}
	if len(entities) != 1 {
		t.Fatalf("1件を期待: got=%d", len(entities))
	}

	// 返却スライスを書き換えてもリポジトリ内部状態に影響しないことを確認する。
	entities = entities[:0]

	again, err := repository.FindAll()
	if err != nil {
		t.Fatalf("再取得に失敗: %v", err)
	}
	if len(again) != 1 {
		t.Fatalf("内部状態が変更されてはならない: got=%d", len(again))
	}
}

func TestDeleteByIDで削除できること(t *testing.T) {
	repository := memory.NewTodoRepository()
	title, _ := domain.NewTitle("牛乳を買う")
	_ = repository.Save(domain.NewEntity("todo-1", title))

	deleted, err := repository.DeleteByID("todo-1")
	if err != nil {
		t.Fatalf("削除に失敗: %v", err)
	}
	if !deleted {
		t.Fatalf("削除結果はtrueであるべき")
	}

	_, ok, _ := repository.FindByID("todo-1")
	if ok {
		t.Fatalf("削除後は見つからないべき")
	}
}

func TestFindByCompletedで完了状態で絞り込めること(t *testing.T) {
	repository := memory.NewTodoRepository()
	title1, _ := domain.NewTitle("牛乳を買う")
	title2, _ := domain.NewTitle("散歩する")

	completed := domain.NewEntity("todo-1", title1)
	completed.Complete()
	_ = repository.Save(completed)
	_ = repository.Save(domain.NewEntity("todo-2", title2))

	entities, err := repository.FindByCompleted(true)
	if err != nil {
		t.Fatalf("絞り込み取得に失敗: %v", err)
	}
	if len(entities) != 1 {
		t.Fatalf("1件を期待: got=%d", len(entities))
	}
	if entities[0].ID() != "todo-1" {
		t.Fatalf("完了Todoのみ返るべき: got=%s", entities[0].ID())
	}
}
