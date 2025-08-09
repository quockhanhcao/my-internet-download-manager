package logic

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewAccountHandler,
	NewHashHandler,
    NewTokenHandler,
)
