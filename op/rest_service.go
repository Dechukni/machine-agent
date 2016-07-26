package op

import "github.com/evoevodin/machine-agent/rest"

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
