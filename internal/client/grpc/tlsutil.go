package grpc

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/ArtShib/gophkeeper/internal/models"
	"google.golang.org/grpc/credentials"
)

func LoadMTLSClient(tlsConfig *models.ConfigTLS) (credentials.TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair(tlsConfig.Cert, tlsConfig.Key)
	if err != nil {
		return nil, err
	}

	caCert, err := os.ReadFile(tlsConfig.CA)
	if err != nil {
		return nil, err
	}
	caPool := x509.NewCertPool()
	caPool.AppendCertsFromPEM(caCert)

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caPool,
		ServerName:   "localhost",
	}

	return credentials.NewTLS(config), nil
}
