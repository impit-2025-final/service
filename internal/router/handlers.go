package router

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
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

	var err error

	if r.Header.Get("Content-Encoding") == "gzip" {
		r.Body, err = handleGzipBody(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
	}

	var dockerInfo domain.DockerInfo

	if err := json.NewDecoder(r.Body).Decode(&dockerInfo); err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(dockerInfo)
	err = h.useCase.UpdateDockerInfo(r.Context(), dockerInfo)
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

	var err error

	if r.Header.Get("Content-Encoding") == "gzip" {
		r.Body, err = handleGzipBody(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
	}

	var networkTraffic []domain.NetworkTraffic

	if err := json.NewDecoder(r.Body).Decode(&networkTraffic); err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(networkTraffic)

	err = h.useCase.UpdateNetworkTraffic(r.Context(), networkTraffic)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleGzipBody(body io.ReadCloser) (io.ReadCloser, error) {
	reader, err := gzip.NewReader(body)
	if err != nil {
		return nil, fmt.Errorf("error: %w", err)
	}
	return reader, nil
}
