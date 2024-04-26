package demo

import (
	"context"
	"errors"
)

type Z struct {
	Val int
}

func NewZ(ctx context.Context, y Y) (Z, error) {
	if y.Val == 0 {
		return Z{}, errors.New("cannot provide z when value is zero")
	}
	return Z{Val: y.Val + 2}, nil
}
