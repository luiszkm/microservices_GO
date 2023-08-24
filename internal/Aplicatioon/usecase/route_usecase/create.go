package routeusecase

import "github.com/luiszkm/microservices_GO/internal/Domain/freight/entity"

type CreateRouteInput struct {
	ID       string
	Name     string
	Distance float64
	Event    string
}
type CreateRouteOutput struct {
	ID           string
	Name         string
	Distance     float64
	Status       string
	FreightPrice float64
}

type CreateRouteUseCase struct {
	Repository entity.RouterRepository
	Freight    entity.FreightInterface
}

func NewCreateRouteUseCase(repository entity.RouterRepository, freight entity.FreightInterface) *CreateRouteUseCase {
	return &CreateRouteUseCase{
		Repository: repository,
		Freight:    freight,
	}
}

func (u *CreateRouteUseCase) Execute(input CreateRouteInput) (*CreateRouteOutput, error) {
	route := entity.NewRoute(input.ID, input.Name, input.Distance)
	u.Freight.Calculate(route)
	err := u.Repository.Create(route)
	if err != nil {
		return nil, err
	}

	return &CreateRouteOutput{
		ID:           route.ID,
		Name:         route.Name,
		Distance:     route.Distance,
		Status:       route.Status,
		FreightPrice: route.FreightPrice,
	}, nil
}
