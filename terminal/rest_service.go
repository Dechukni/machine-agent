package terminal

import (
	"github.com/evoevodin/machine-agent/rest"
	"net/http"
)

var HttpRoutes = rest.RoutesGroup{
	"Terminal routes",
	[]rest.Route{
		{
			"GET",
			"Connect to pty(webscoket)",
			"/pty",
			ConnectToPtyHF,
		},
	},
}

func ConnectToPtyHF(w http.ResponseWriter, r *http.Request) error {
	ptyHandler(w, r)
	return nil
}
