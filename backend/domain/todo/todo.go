package todo

import (
	"errors"
	"strings"
)

const maxTitleLength = 80

var (
	ErrTitleRequired  = errors.New("タイトルは必須です")
	ErrTitleTooLong   = errors.New("タイトルは80文字以内で入力してください")
	ErrIDRequired     = errors.New("IDは必須です")
	ErrTitleIsInvalid = errors.New("タイトルが不正です")
)

// TitleはTodoのタイトルを表す値オブジェクト。
// ドメインの制約を生成時に必ず満たすことで、以後の処理を単純化する。
type Title struct {
	value string
}

func NewTitle(value string) (Title, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return Title{}, ErrTitleRequired
	}
	if len([]rune(trimmed)) > maxTitleLength {
		return Title{}, ErrTitleTooLong
	}
	return Title{value: trimmed}, nil
}

func (t Title) Value() string {
	return t.value
}

// EntityはTodoのエンティティ。
// 生成時点で不正状態を持たせないことを前提にする。
type Entity struct {
	id          string
	title       Title
	isCompleted bool
}

func NewEntity(id string, title Title) Entity {
	if strings.TrimSpace(id) == "" {
		panic(ErrIDRequired)
	}
	if title.value == "" {
		panic(ErrTitleIsInvalid)
	}
	return Entity{id: id, title: title, isCompleted: false}
}

func (e Entity) ID() string {
	return e.id
}

func (e Entity) Title() Title {
	return e.title
}

func (e Entity) IsCompleted() bool {
	return e.isCompleted
}

// CompleteはTodoを完了状態に遷移させる。
// 完了済みに再適用しても状態を壊さないよう、冪等に扱う。
func (e *Entity) Complete() {
	e.isCompleted = true
}

// ChangeTitleはTodoのタイトルを差し替える。
// タイトルの妥当性は値オブジェクト生成時に保証される前提で受け取る。
func (e *Entity) ChangeTitle(title Title) {
	e.title = title
}
