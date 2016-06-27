package op

import "github.com/evoevodin/machine-agent/core"

var (
	HttpRoutes = core.HttpRoutesGroup{
		"Dispatcher Http Routes",
		[]core.HttpRoute{
			core.HttpRoute{
				"GET",
				"Connect to Machine-Agent",
				"/connect",
				RegisterChannel,
			},
		},
	}
)
