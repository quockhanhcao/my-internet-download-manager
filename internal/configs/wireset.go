package configs

import "github.com/google/wire"

var WireSet = wire.NewSet(
	NewConfig,
	wire.FieldsOf(new(Config), "DatabaseConfig"),
	wire.FieldsOf(new(Config), "AuthConfig"),
	wire.FieldsOf(new(Config), "LogConfig"),
)
