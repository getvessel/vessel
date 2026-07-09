package terminal

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gorilla/websocket"
	"vessel.dev/vessel/internal/middleware"
	"vessel.dev/vessel/internal/utils"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type AppService struct {
	ID          string
	ContainerID string
}

type ServiceRepository interface {
	GetByID(ctx context.Context, id string) (*AppService, error)
}

type Handler struct {
	dockerClient  *client.Client
	tokenService  TokenValidator
	serviceRepo   ServiceRepository
	normalizeName func(id string) string
}

type TokenValidator interface {
	ValidateToken(tokenStr string) (*TokenClaim, error)
}

type TokenClaim struct {
	UserID string
	Email  string
}

func NewHandler(
	dockerClient *client.Client,
	tokenService TokenValidator,
	serviceRepo ServiceRepository,
) *Handler {
	return &Handler{
		dockerClient:  dockerClient,
		tokenService:  tokenService,
		serviceRepo:   serviceRepo,
		normalizeName: utils.NormalizeContainerName,
	}
}

func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	if h.tokenService != nil {
		tokenStr := middleware.ExtractTokenFromRequest(r)
		if tokenStr == "" {
			writeError(w, http.StatusUnauthorized, "missing authentication token for terminal access")
			return
		}
		if _, err := h.tokenService.ValidateToken(tokenStr); err != nil {
			writeError(w, http.StatusUnauthorized, "invalid authentication token for terminal access")
			return
		}
	}

	id := r.PathValue("id")
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing id parameter")
		return
	}

	containerName := h.normalizeName(id)
	if h.serviceRepo != nil {
		if svc, err := h.serviceRepo.GetByID(r.Context(), id); err == nil && svc != nil {
			if svc.ContainerID != "" && svc.ContainerID != "-" {
				containerName = svc.ContainerID
			} else {
				containerName = h.normalizeName(svc.ID)
			}
		}
	}

	execConfig := types.ExecConfig{
		Cmd:          []string{"/bin/sh"},
		Tty:          true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
	}

	if h.dockerClient == nil {
		writeError(w, http.StatusServiceUnavailable, "docker client unavailable")
		return
	}

	ctx := r.Context()
	execID, err := h.dockerClient.ContainerExecCreate(ctx, containerName, execConfig)
	if err != nil {
		execConfig.Cmd = []string{"/bin/sh"}
		execID, err = h.dockerClient.ContainerExecCreate(ctx, containerName, execConfig)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "failed to create interactive container session: "+err.Error())
			return
		}
	}

	execResp, err := h.dockerClient.ContainerExecAttach(ctx, execID.ID, types.ExecStartCheck{Tty: true})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to attach to container shell: "+err.Error())
		return
	}
	defer execResp.Close()

	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer wsConn.Close()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		defer cancel()
		buf := make([]byte, 4096)
		for {
			n, err := execResp.Reader.Read(buf)
			if err != nil {
				if err != io.EOF {
					_ = wsConn.WriteMessage(websocket.TextMessage, []byte("\r\n[session terminated]\r\n"))
				}
				break
			}
			if n > 0 {
				_ = wsConn.SetWriteDeadline(time.Now().Add(10 * time.Second))
				if err := wsConn.WriteMessage(websocket.TextMessage, buf[:n]); err != nil {
					break
				}
			}
		}
	}()

	for {
		_, message, err := wsConn.ReadMessage()
		if err != nil {
			break
		}
		if len(message) > 0 {
			if _, err := execResp.Conn.Write(message); err != nil {
				break
			}
		}
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	http.Error(w, msg, status)
}
