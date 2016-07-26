package op

import (
	"github.com/evoevodin/machine-agent/core/rest"
)

var HttpRoutes = rest.RoutesGroup{
	"Channel Routes",
	[]rest.Route{
		{
			"GET",
			"Connect to Machine-Agent(webscoket)",
			"/connect",
			registerChannel,
		},
	},
}
