package http

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
)

type config struct {
	tlsSkipVerify bool
	certificates  []string
}

// Option is a configuration option for the HTTP client.
type Option func(c *config)

// WithTLSSkipVerify enables or disables the verification of the server TLS
// certificates. This option defaults to false.
func WithTLSSkipVerify(tlsSkipVerify bool) Option {
	return func(c *config) {
		c.tlsSkipVerify = tlsSkipVerify
	}
}

// WithExtraCertificates adds more certificates to the pool of certificates
// trusted by this client. path is the path to a PEM file containing one or more
// certificates. If path is empty or the file at path does not exist, no
// certificates are added to the pool.
func WithExtraCertificates(path string) Option {
	return func(c *config) {
		c.certificates = append(c.certificates, path)
	}
}

func NewClient(options ...Option) (*http.Client, error) {
	var c config

	for _, o := range options {
		o(&c)
	}

	pool, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf("read system cert pool: %w", err)
	}
	if pool == nil {
		pool = x509.NewCertPool()
	}

	if err := loadCertificates(pool, c.certificates); err != nil {
		return nil, fmt.Errorf("load certificates: %w", err)
	}

	client := http.DefaultClient

	if transport, ok := http.DefaultTransport.(*http.Transport); ok {
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: c.tlsSkipVerify,
			RootCAs:            pool,
		}
	}

	return client, nil
}

func loadCertificates(pool *x509.CertPool, certificates []string) error {
	for _, certificate := range certificates {
		if err := loadCertificate(pool, certificate); err != nil {
			return fmt.Errorf("load certificate %v: %w", certificate, err)
		}
	}

	return nil
}

func loadCertificate(pool *x509.CertPool, certPath string) error {
	certData, err := os.ReadFile(certPath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read certificates: %w", err)
	}

	if ok := pool.AppendCertsFromPEM(certData); !ok {
		return fmt.Errorf("no certificates found from NODE_EXTRA_CA_CERTS")
	}

	return nil
}
