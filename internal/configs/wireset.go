package configs

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewConfig,
	wire.FieldsOf(new(Config), "AccountConfig"),
	wire.FieldsOf(new(Config), "DatabaseConfig"),
)
