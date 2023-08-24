package routeusecase

import (
	"time"

	"github.com/luiszkm/microservices_GO/internal/Domain/freight/entity"
)

type ChangeStatusInput struct {
	ID         string
	StartedAt  entity.CustomTime
	FinishedAt entity.CustomTime
	Event      string
}

type ChangeStatusOutput struct {
	ID         string
	Status     string
	StartedAt  entity.CustomTime
	FinishedAt entity.CustomTime
}

type ChangeStatusUseCase struct {
	Repository entity.RouterRepository
}

func NewChangeRouteStatusUseCase(repository entity.RouterRepository) *ChangeStatusUseCase {
	return &ChangeStatusUseCase{
		Repository: repository,
	}
}

func (u *ChangeStatusUseCase) Execute(input ChangeStatusInput) (*ChangeStatusOutput, error) {
	route, err := u.Repository.FindById(input.ID)
	if err != nil {
		return nil, err
	}

	switch input.Event {
	case "start":
		route.Start(time.Time(input.StartedAt))
	case "finish":
		route.Finish(time.Time(input.FinishedAt))
	}

	err = u.Repository.Update(route)
	if err != nil {
		return nil, err
	}

	return &ChangeStatusOutput{
		ID:         route.ID,
		Status:     route.Status,
		StartedAt:  entity.CustomTime(route.StartedAt),
		FinishedAt: entity.CustomTime(route.FinishedAt),
	}, nil
}