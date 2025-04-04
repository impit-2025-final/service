package router

import (
	"net/http"
)

func SetupRoutes(handler *Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/docker-info", handler.CreateDockerInfo)
	mux.HandleFunc("/network-traffic", handler.CreateNetworkTraffic)

	mux.HandleFunc("/node-info", handler.CreateNodeInfo)
	mux.HandleFunc("/node-info-update", handler.UpdateNodeInfo)

	return mux
}
