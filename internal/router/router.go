package router

import (
	"net/http"
)

func SetupRoutes(handler *Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/docker-info", handler.GetDockerInfo)
	mux.HandleFunc("/network-traffic", handler.GetNetworkTraffic)

	return mux
}
