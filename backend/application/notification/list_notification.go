package notification

import domain "github.com/user/ddd/backend/domain/notification"

type ListNotificationUseCase struct {
	repository Repository
}

func NewListNotificationUseCase(repository Repository) ListNotificationUseCase {
	return ListNotificationUseCase{repository: repository}
}

func (u ListNotificationUseCase) Execute() ([]domain.Notification, error) {
	return u.repository.FindAll()
}
