package op

import (
	"github.com/evoevodin/machine-agent/core/rest"
)

var HttpRoutes = rest.HttpRoutesGroup{
	"Dispatcher Http Routes",
	[]rest.HttpRoute{
		rest.HttpRoute{
			"GET",
			"Connect to Machine-Agent(webscoket)",
			"/connect",
			RegisterChannel,
		},
	},
}
