package handler

import (
	"github.com/google/wire"
	"github.com/quockhanhcao/my-internet-download-manager/internal/handler/grpc"
	"github.com/quockhanhcao/my-internet-download-manager/internal/handler/http"
)

var WireSet = wire.NewSet(
    grpc.WireSet,
    http.WireSet,
)
