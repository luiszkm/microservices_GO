package entity

import "time"

type CustomTime time.Time
const layout = "2006-01-02T15:04:05.000Z"

func (ct *CustomTime) UnmarshalJSON(b []byte) (err error) {
	t, err := time.Parse(layout, string(b))
	if err != nil {
		return err
	}
	*ct = CustomTime(t)
	return nil
}

type RouterRepository interface {
	Create(route *Route) error
	FindById(id string) (*Route, error)
	Update(route *Route) error
}

type Route struct {
	ID           string
	Name         string
	Distance     float64
	Status       string
	FreightPrice float64
	StartedAt    time.Time
	FinishedAt   time.Time
}

func (r *Route) Start(startedAt time.Time) {
	r.Status = "started"
	r.StartedAt = startedAt
}
func (r *Route) Finish(finishedAt time.Time) {
	r.Status = "finished"
	r.FinishedAt = finishedAt
}


func NewRoute (id, name string, distance float64) *Route {
	return &Route{
		ID: id,
		Name: name,
		Distance: distance,
		Status: "pending",
	}
}