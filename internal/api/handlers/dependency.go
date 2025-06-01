package handlers

import (
	"canary/internal/services/dependency"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"

	"canary/internal/logger"
	"canary/internal/types"
)

type DependencyHandler struct {
	service dependency.Service
}

func NewDependencyHandler(service dependency.Service) *DependencyHandler {
	return &DependencyHandler{
		service: service,
	}
}

func (h *DependencyHandler) ResolveDependency() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DependencyCheckRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			RespondJSON(w, http.StatusBadRequest, map[string]string{
				"status":  "error",
				"message": "Invalid request body",
			})
			return
		}

		if req.Ecosystem == "" || req.Package == "" || req.CurrentVersion == "" || req.UpgradeVersion == "" || req.FileContent == "" {
			RespondJSON(w, http.StatusBadRequest, map[string]string{
				"status":  "error",
				"message": "Missing required fields. Required: ecosystem, package, currentVersion, upgradeVersion, fileContent",
			})
			return
		}

		updatedContent, err := h.service.Resolve(&req)
		if err != nil {
			status := http.StatusInternalServerError
			errMsg := "Failed to process dependency update"
			if errors.Is(err, base64.CorruptInputError(0)) {
				status = http.StatusBadRequest
				errMsg = "Invalid base64 file content"
			}

			logger.Error("error resolving dependency", "error", err)
			RespondJSON(w, status, map[string]string{
				"status":  "error",
				"message": errMsg,
			})
			return
		}

		response := types.DependencyCheckResponse{
			FileContent: updatedContent,
		}

		RespondJSON(w, http.StatusOK, response)
	}
}
