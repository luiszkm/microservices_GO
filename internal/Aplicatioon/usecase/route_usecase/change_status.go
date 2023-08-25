package routeusecase

import (
	"time"

	"github.com/luiszkm/microservices_GO/internal/Domain/freight/entity"
)

type ChangeStatusInput struct {
	ID         string            `json:"id"`
	StartedAt  entity.CustomTime `json:"started_at"`
	FinishedAt entity.CustomTime `json:"finished_at"`
	Event      string            `json:"event"`
}

type ChangeStatusOutput struct {
	ID         string            `json:"id"`
	Status     string            `json:"status"`
	StartedAt  entity.CustomTime `json:"started_at"`
	FinishedAt entity.CustomTime `json:"finished_at"`
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
	case "RouteStarted":
		route.Start(time.Time(input.StartedAt))
	case "RouteFinished":
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
