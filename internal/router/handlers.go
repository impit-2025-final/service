package router

import (
	"compress/gzip"
	"crypto/rand"
	"encoding/hex"
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

func (h *Handler) CreateDockerInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nodeId, err := h.checkToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

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
	err = h.useCase.UpdateDockerInfo(r.Context(), dockerInfo, nodeId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) CreateNetworkTraffic(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	_, err := h.checkToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if r.Header.Get("Content-Encoding") == "gzip" {
		r.Body, err = handleGzipBody(r.Body)
		if err != nil {
			fmt.Println(err)
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

func (h *Handler) CreateNodeInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var nodeInfo domain.NodeInfo

	if err := json.NewDecoder(r.Body).Decode(&nodeInfo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token := GenerateSecureToken(32)

	nodeInfo, err := h.useCase.CreateNodeInfo(r.Context(), string(token), nodeInfo.NodeName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(nodeInfo)
}

func (h *Handler) UpdateNodeInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var nodeInfo domain.NodeInfo
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	oldNodeInfo, err := h.useCase.GetNodeInfo(r.Context(), token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&nodeInfo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	oldNodeInfo.Hostname = nodeInfo.Hostname
	oldNodeInfo.Ips = nodeInfo.Ips

	fmt.Println(oldNodeInfo)

	err = h.useCase.UpdateNodeInfo(r.Context(), oldNodeInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) checkToken(r *http.Request) (uint, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return 0, fmt.Errorf("unauthorized")
	}

	nodeInfo, err := h.useCase.GetNodeInfo(r.Context(), token)
	if err != nil {
		return 0, fmt.Errorf("unauthorized")
	}
	return nodeInfo.ID, nil
}

func handleGzipBody(body io.ReadCloser) (io.ReadCloser, error) {
	reader, err := gzip.NewReader(body)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("error: %w", err)
	}
	return reader, nil
}

func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
