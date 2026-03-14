package notification

import (
	"fmt"
	"strings"
	"time"

	domain "github.com/user/ddd/backend/domain/notification"
)

// Repositoryは通知コンテキストの永続化抽象。
type Repository interface {
	Save(notification domain.Notification) error
	FindAll() ([]domain.Notification, error)
}

// IDGeneratorは通知ID生成の抽象。
type IDGenerator interface {
	NextID() string
}

// Clockは時刻取得の抽象。
type Clock interface {
	Now() time.Time
}

type RecordTodoCompletedCommand struct {
	TodoID string
	Title  string
}

type RecordTodoCompletedUseCase struct {
	repository  Repository
	idGenerator IDGenerator
	clock       Clock
}

func NewRecordTodoCompletedUseCase(
	repository Repository,
	idGenerator IDGenerator,
	clock Clock,
) RecordTodoCompletedUseCase {
	return RecordTodoCompletedUseCase{repository: repository, idGenerator: idGenerator, clock: clock}
}

func (u RecordTodoCompletedUseCase) Execute(command RecordTodoCompletedCommand) (domain.Notification, error) {
	message := fmt.Sprintf("Todo[%s]を完了しました: %s", strings.TrimSpace(command.TodoID), strings.TrimSpace(command.Title))
	notification, err := domain.NewNotification(u.idGenerator.NextID(), message, u.clock.Now())
	if err != nil {
		return domain.Notification{}, err
	}
	if err := u.repository.Save(notification); err != nil {
		return domain.Notification{}, err
	}
	return notification, nil
}
