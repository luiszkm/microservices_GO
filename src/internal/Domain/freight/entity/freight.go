package entity

type FreightInterface interface {
	Calculate(route *Route)
}

type Freight struct {
	PricePerKm float64
}

func (f *Freight) Calculate(route *Route)  {
	  route.FreightPrice = route.Distance * f.PricePerKm
}
func NewFreight(pricePerKm float64) *Freight {
	return &Freight{
		PricePerKm: pricePerKm,
	}
}