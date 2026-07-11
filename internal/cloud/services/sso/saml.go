package sso

import (
	"crypto/rsa"
	"crypto/x509"
	"io"
	"net/http"
	"net/url"

	"github.com/crewjam/saml/samlsp"
)

type SAMLService struct {
	BaseURL string
	Key     *rsa.PrivateKey
	Cert    *x509.Certificate
}

func NewSAMLService(baseURL string, key *rsa.PrivateKey, cert *x509.Certificate) *SAMLService {
	return &SAMLService{
		BaseURL: baseURL,
		Key:     key,
		Cert:    cert,
	}
}

func (s *SAMLService) ConfigureMiddleware(idpMetadataURL string) (*samlsp.Middleware, error) {
	idpMetadata, err := url.Parse(idpMetadataURL)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(idpMetadata.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	metadataBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	idpEntityDescriptor, err := samlsp.ParseMetadata(metadataBytes)
	if err != nil {
		return nil, err
	}

	rootURL, err := url.Parse(s.BaseURL)
	if err != nil {
		return nil, err
	}

	return samlsp.New(samlsp.Options{
		URL:         *rootURL,
		Key:         s.Key,
		Certificate: s.Cert,
		IDPMetadata: idpEntityDescriptor,
	})
}
