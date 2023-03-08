/*
 * © 2023 Snyk Limited
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package http

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestWithTLSSkipVerify(t *testing.T) {
	privateKey := createPrivateKey(t)
	certificate := createSelfSignedCertificate(t, privateKey)

	certificatePEM := encodeCertificateToFile(t, certificate)
	privateKeyPEM := encodePrivateKeyToFile(t, privateKey)

	serverAddr := startHTTPSServer(t, certificatePEM, privateKeyPEM)

	// When skipping the TLS verification, the client doesn't verify the server
	// certificates. Requests will not fail even if the server uses a self-signed
	// certificate.

	client, err := NewClient(WithTLSSkipVerify(true))
	if err != nil {
		t.Fatalf("create client: %v", err)
	}

	req, err := http.NewRequest(http.MethodGet, serverAddr, nil)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}

	if _, err := client.Do(req); err != nil {
		t.Fatalf("perform request: %v", err)
	}
}

func TestWithExtraCertificates(t *testing.T) {
	privateKey := createPrivateKey(t)
	certificate := createSelfSignedCertificate(t, privateKey)

	certificatePEM := encodeCertificateToFile(t, certificate)
	privateKeyPEM := encodePrivateKeyToFile(t, privateKey)

	serverAddr := startHTTPSServer(t, certificatePEM, privateKeyPEM)

	// When using extra certificates, the client will accept as valid whatever TLS
	// certificate it has in its own pool. If the pool includes the self-signed
	// certificate from the server, requests will not fail.

	client, err := NewClient(WithExtraCertificates(certificatePEM))
	if err != nil {
		t.Fatalf("create client: %v", err)
	}

	req, err := http.NewRequest(http.MethodGet, serverAddr, nil)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}

	if _, err := client.Do(req); err != nil {
		t.Fatalf("perform request: %v", err)
	}
}

func startHTTPSServer(t *testing.T, certificatePEM, privateKeyPEM string) string {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := http.Server{
		Handler: handler,
	}

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)

		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			t.Fatalf("shutdown server: %v", err)
		}
	})

	go func() {
		if err := server.ServeTLS(listener, certificatePEM, privateKeyPEM); err != http.ErrServerClosed {
			t.Logf("serve TLS: %v", err)
		}
	}()

	return fmt.Sprintf("https://%s", listener.Addr())
}

func createPrivateKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("create private key: %v", err)
	}

	return privateKey
}

func createSelfSignedCertificate(t *testing.T, privateKey *rsa.PrivateKey) []byte {
	t.Helper()

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Organization"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 180),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
	}

	certificate, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatalf("create certificate: %v", err)
	}

	return certificate
}

func encodeCertificateToFile(t *testing.T, certificate []byte) string {
	file, err := os.Create(filepath.Join(t.TempDir(), "cert.pem"))
	if err != nil {
		t.Fatalf("create file: %v", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			t.Fatalf("close file: %v", err)
		}
	}()

	encodeCertificate(t, file, certificate)

	return file.Name()
}

func encodePrivateKeyToFile(t *testing.T, key *rsa.PrivateKey) string {
	file, err := os.Create(filepath.Join(t.TempDir(), "key.pem"))
	if err != nil {
		t.Fatalf("create file: %v", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			t.Fatalf("close file: %v", err)
		}
	}()

	encodePrivateKey(t, file, key)

	return file.Name()
}

func encodeCertificate(t *testing.T, w io.Writer, certificate []byte) {
	t.Helper()

	if err := pem.Encode(w, &pem.Block{Type: "CERTIFICATE", Bytes: certificate}); err != nil {
		t.Fatalf("encode certificate: %v", err)
	}
}

func encodePrivateKey(t *testing.T, w io.Writer, key *rsa.PrivateKey) {
	t.Helper()

	if err := pem.Encode(w, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}); err != nil {
		t.Fatalf("encode private key: %v", err)
	}
}
