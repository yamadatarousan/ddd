package todo

import (
	"errors"
	"strings"

	domain "github.com/user/ddd/backend/domain/todo"
)

var ErrIDGenerationFailed = errors.New("IDの生成に失敗しました")

// RepositoryはTodo永続化の抽象。
// ユースケースは永続化方式を知らないように依存方向を固定する。
type Repository interface {
	Save(entity domain.Entity) error
}

type CreateTodoCommand struct {
	Title string
}

type CreateTodoUseCase struct {
	repository Repository
	generateID func() string
}

func NewCreateTodoUseCase(repository Repository, generateID func() string) CreateTodoUseCase {
	return CreateTodoUseCase{repository: repository, generateID: generateID}
}

func (u CreateTodoUseCase) Execute(command CreateTodoCommand) (domain.Entity, error) {
	title, err := domain.NewTitle(command.Title)
	if err != nil {
		return domain.Entity{}, err
	}

	id := strings.TrimSpace(u.generateID())
	if id == "" {
		return domain.Entity{}, ErrIDGenerationFailed
	}

	entity := domain.NewEntity(id, title)
	if err := u.repository.Save(entity); err != nil {
		return domain.Entity{}, err
	}

	return entity, nil
}
