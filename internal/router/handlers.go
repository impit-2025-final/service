package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	domain "service/internal/domain"
	"service/internal/usecase"
)

type Handler struct {
	useCase usecase.ContainerUseCase
}

func NewHandler(useCase usecase.ContainerUseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) GetDockerInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var dockerInfo domain.DockerInfo

	if err := json.NewDecoder(r.Body).Decode(&dockerInfo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.useCase.UpdateDockerInfo(r.Context(), dockerInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) GetNetworkTraffic(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var networkTraffic []domain.NetworkTraffic

	if err := json.NewDecoder(r.Body).Decode(&networkTraffic); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(networkTraffic)

	err := h.useCase.UpdateNetworkTraffic(r.Context(), networkTraffic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
