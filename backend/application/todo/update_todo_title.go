package todo

import (
	"strings"

	domain "github.com/user/ddd/backend/domain/todo"
)

// UpdateTitleRepositoryはタイトル変更ユースケースが必要とする永続化操作。
// 対象取得と保存のみを定義し、ユースケース外の責務を持ち込まない。
type UpdateTitleRepository interface {
	FindByID(id string) (domain.Entity, bool, error)
	Save(entity domain.Entity) error
}

type UpdateTodoTitleCommand struct {
	ID    string
	Title string
}

type UpdateTodoTitleUseCase struct {
	repository UpdateTitleRepository
}

func NewUpdateTodoTitleUseCase(repository UpdateTitleRepository) UpdateTodoTitleUseCase {
	return UpdateTodoTitleUseCase{repository: repository}
}

func (u UpdateTodoTitleUseCase) Execute(command UpdateTodoTitleCommand) (domain.Entity, error) {
	id := strings.TrimSpace(command.ID)
	if id == "" {
		return domain.Entity{}, ErrTodoIDRequired
	}

	entity, ok, err := u.repository.FindByID(id)
	if err != nil {
		return domain.Entity{}, err
	}
	if !ok {
		return domain.Entity{}, ErrTodoNotFound
	}

	title, err := domain.NewTitle(command.Title)
	if err != nil {
		return domain.Entity{}, err
	}

	entity.ChangeTitle(title)
	if err := u.repository.Save(entity); err != nil {
		return domain.Entity{}, err
	}

	return entity, nil
}
