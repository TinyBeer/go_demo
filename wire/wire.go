//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"learn_wire/demo"

	"github.com/google/wire"
)

func InitializeZ(cxt context.Context) (demo.Z, error) {
	wire.Build(demo.ProviderSet)
	return demo.Z{}, nil
}
