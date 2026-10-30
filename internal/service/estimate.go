// Package service contains illustrative domain logic, independent of HTTP.
package service

import (
	"context"
	"errors"
)

type Offering struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	BaseUnits int    `json:"base_units"`
}

// These examples are not Cynbit's products, prices or effort estimates.
func Catalog() []Offering {
	return []Offering{{"api-review", "Illustrative API review", 3}, {"integration-review", "Illustrative integration review", 5}}
}

type Input struct {
	ServiceID    string `json:"service_id"`
	Integrations int    `json:"integrations"`
}

type Estimate struct {
	ServiceID   string `json:"service_id"`
	EffortUnits int    `json:"effort_units"`
	Disclaimer  string `json:"disclaimer"`
}

var ErrInvalid = errors.New("invalid estimate input")

type Estimator interface {
	Estimate(context.Context, Input) (Estimate, error)
}
type Calculator struct{}

func (Calculator) Estimate(ctx context.Context, in Input) (Estimate, error) {
	if err := ctx.Err(); err != nil {
		return Estimate{}, err
	}
	if in.Integrations < 0 || in.Integrations > 20 {
		return Estimate{}, ErrInvalid
	}
	for _, o := range Catalog() {
		if o.ID == in.ServiceID {
			return Estimate{in.ServiceID, o.BaseUnits + 2*in.Integrations, "Illustrative units only; not a quote or Cynbit service."}, nil
		}
	}
	return Estimate{}, ErrInvalid
}
