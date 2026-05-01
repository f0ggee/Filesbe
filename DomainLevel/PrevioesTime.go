package DomainLevel

import (
	"errors"
	"time"
)

const (
	ValueError = "variables cannot be zero or a default value"
)

type PreviousSwapTime struct {
	Id   string
	Time time.Time
}

func (p *PreviousSwapTime) NewPreviousTime(id string, t time.Time) error {
	if len(id) == 0 || t.IsZero() {

		return errors.New(ValueError)
	}
	p.Id = id
	p.Time = t
	return nil
}
func (this PreviousSwapTime) GetPreviousSwapTime() time.Time { return this.Time }

func Get() *PreviousSwapTime {

	return &PreviousSwapTime{}
}
