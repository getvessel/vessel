package api

import (
	"fmt"
	"net/http"
	"path/filepath"

	"vessel.dev/vessel/internal/types"
)

func (s *Server) handleDeployProject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing project id parameter")
		return
	}

	project, err := s.store.GetProject(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "project not found")
		return
	}

	sourceDir := filepath.Join("data", "builds", id)
	if s.gitService != nil {
		_ = s.gitService.CloneOrPullRepository(r.Context(), id, sourceDir, nil)
	}

	containerID, err := s.deployer.Deploy(r.Context(), project, sourceDir, nil)
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("deployment rollout failed: %v", err))
		return
	}

	_ = s.proxyManager.Reload(r.Context())
	writeJSON(w, http.StatusOK, map[string]string{
		"status":       "deployed",
		"container_id": containerID,
	})
}
