package demo

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewX,
	NewY,
	NewZ,
)
