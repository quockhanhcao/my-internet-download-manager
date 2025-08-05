//go:build wireinject
// +build wireinject

//
//go:generate go run github.com/google/wire/cmd/wire
package wiring

import (
	"github.com/google/wire"
	"github.com/quockhanhcao/my-internet-download-manager/internal/configs"
	"github.com/quockhanhcao/my-internet-download-manager/internal/dataacess"
	"github.com/quockhanhcao/my-internet-download-manager/internal/handler"
	"github.com/quockhanhcao/my-internet-download-manager/internal/handler/grpc"
	"github.com/quockhanhcao/my-internet-download-manager/internal/logic"
)

var WireSet = wire.NewSet(
	configs.WireSet,
	dataacess.WireSet,
	logic.WireSet,
	handler.WireSet,
)

func InitializeGRPCServer(configFilePath configs.ConfigFilePath) (grpc.Server, func(), error) {
	wire.Build(WireSet)
	return nil, nil, nil
}
