package dataacess

import (
	"github.com/google/wire"
	"github.com/quockhanhcao/my-internet-download-manager/internal/dataacess/database"
)

var WireSet = wire.NewSet(
	database.WireSet,
)
