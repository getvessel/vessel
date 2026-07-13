package handlers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"time"

	"vessl.dev/vessl/internal/cloud/services"

	"github.com/crewjam/saml/samlsp"
	"github.com/labstack/echo/v4"
)

type SSOHandler struct {
	samlMiddleware *samlsp.Middleware
}

func NewSSOHandler() (*SSOHandler, error) {
	metadataURL := os.Getenv("SAML_IDP_METADATA_URL")
	if metadataURL == "" {
		return nil, fmt.Errorf("SAML_IDP_METADATA_URL not set")
	}

	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	template := x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{Organization: []string{"Vessl Cloud"}},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 365),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	certBytes, _ := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	cert, _ := x509.ParseCertificate(certBytes)

	baseURL := os.Getenv("VESSL_CLOUD_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	svc := services.NewSAMLService(baseURL, key, cert)
	sp, err := svc.ConfigureMiddleware(metadataURL)
	if err != nil {
		return nil, fmt.Errorf("failed to configure SAML middleware: %v", err)
	}
	return &SSOHandler{samlMiddleware: sp}, nil
}

func (h *SSOHandler) RegisterRoutes(g *echo.Group) {
	g.Any("/saml/*", echo.WrapHandler(h.samlMiddleware))
}

func (h *SSOHandler) RequireSAML() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			handler := h.samlMiddleware.RequireAccount(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				c.SetRequest(r)
				if session := samlsp.SessionFromContext(r.Context()); session != nil {
					c.Set("saml_session", session)
				}
				next(c)
			}))
			handler.ServeHTTP(c.Response().Writer, c.Request())
			return nil
		}
	}
}
