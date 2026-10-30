package service

import (
	"context"
	"errors"
	"testing"
)

func TestCalculatorCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := (Calculator{}).Estimate(ctx, Input{"api-review", 1})
	if !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}
func TestCatalogIsolation(t *testing.T) {
	c := Catalog()
	c[0].ID = "mutated"
	if Catalog()[0].ID != "api-review" {
		t.Fatal("catalog mutation escaped")
	}
}
