package handlers

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/docker/docker/client"
	"github.com/gorilla/websocket"
	"github.com/hashicorp/yamux"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, restrict this
	},
}

type AgentHandler struct {
	// A map to hold active connections by token or tenant ID
	activeAgents map[string]*yamux.Session
	mu           sync.RWMutex
}

func NewAgentHandler() *AgentHandler {
	return &AgentHandler{
		activeAgents: make(map[string]*yamux.Session),
	}
}

// AcceptConnection handles incoming websocket connections from remote agents
func (h *AgentHandler) AcceptConnection(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	if token == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Missing Authorization header"})
	}

	// TODO: Validate token against PostgreSQL (CloudDB) to get tenant/server ID
	// For now, we just use the raw token as the identifier
	serverID := token

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		log.Printf("Failed to upgrade agent connection: %v", err)
		return err
	}

	netConn := &websocketConn{conn: ws}
	session, err := yamux.Client(netConn, yamux.DefaultConfig())
	if err != nil {
		ws.Close()
		log.Printf("Failed to establish yamux client session: %v", err)
		return err
	}

	h.mu.Lock()
	h.activeAgents[serverID] = session
	h.mu.Unlock()

	log.Printf("Agent connected securely for server/tenant: %s", serverID)

	// Keep the connection alive until the agent disconnects
	<-session.CloseChan()

	h.mu.Lock()
	delete(h.activeAgents, serverID)
	h.mu.Unlock()

	log.Printf("Agent disconnected: %s", serverID)
	return nil
}

type FleetDeployRequest struct {
	ImageTag    string   `json:"image_tag"`
	AgentTokens []string `json:"agent_tokens"` // Target servers
	EnvVars     []string `json:"env_vars"`
}

// DeployToFleet broadcasts a deployment command to selected agents
// @Summary Fleet Deployment
// @Description Dispatch a deployment instruction to a subset of connected Vossel Daemons
// @Tags Cloud-Fleet
// @Accept json
// @Produce json
// @Success 202 {object} map[string]interface{}
// @Router /cloud/fleet/deploy [post]
func (h *AgentHandler) DeployToFleet(c echo.Context) error {
	var req FleetDeployRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	successfulDispatches := 0
	failedAgents := []string{}

	// In a real implementation, we would encode a JSON instruction
	// and write it to the Yamux session streams.
	for _, token := range req.AgentTokens {
		session, exists := h.activeAgents[token]
		if !exists || session.IsClosed() {
			failedAgents = append(failedAgents, token)
			continue
		}

		// Open a stream to the remote agent
		stream, err := session.Open()
		if err != nil {
			failedAgents = append(failedAgents, token)
			continue
		}

		// Mock payload dispatch
		// e.g. json.NewEncoder(stream).Encode(instruction)
		stream.Write([]byte("DEPLOY:" + req.ImageTag + "\n"))
		stream.Close()

		successfulDispatches++
	}

	return c.JSON(http.StatusAccepted, map[string]interface{}{
		"message":   "Fleet deployment dispatched",
		"successes": successfulDispatches,
		"failures":  failedAgents,
	})
}

// GetDockerClient returns a Docker client that routes traffic over the yamux tunnel
func (h *AgentHandler) GetDockerClient(serverID string) (*client.Client, error) {
	h.mu.RLock()
	session, exists := h.activeAgents[serverID]
	h.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("agent %s is not currently connected", serverID)
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				return session.Open()
			},
		},
	}

	cli, err := client.NewClientWithOpts(
		client.WithHTTPClient(httpClient),
		client.WithHost("http://localhost"), // Host is ignored because of DialContext
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client over tunnel: %w", err)
	}

	return cli, nil
}

// websocketConn wraps a gorilla/websocket to implement net.Conn
type websocketConn struct {
	conn *websocket.Conn
}

func (c *websocketConn) Read(p []byte) (int, error) {
	_, r, err := c.conn.NextReader()
	if err != nil {
		return 0, err
	}
	return r.Read(p)
}

func (c *websocketConn) Write(p []byte) (int, error) {
	err := c.conn.WriteMessage(websocket.BinaryMessage, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (c *websocketConn) Close() error {
	return c.conn.Close()
}

func (c *websocketConn) LocalAddr() net.Addr                { return c.conn.LocalAddr() }
func (c *websocketConn) RemoteAddr() net.Addr               { return c.conn.RemoteAddr() }
func (c *websocketConn) SetDeadline(t time.Time) error      { return c.conn.SetWriteDeadline(t) }
func (c *websocketConn) SetReadDeadline(t time.Time) error  { return c.conn.SetReadDeadline(t) }
func (c *websocketConn) SetWriteDeadline(t time.Time) error { return c.conn.SetWriteDeadline(t) }
