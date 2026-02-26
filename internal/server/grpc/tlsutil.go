package grpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/ArtShib/gophkeeper/internal/models"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
)

func LoadMTLSServer(tlsConfig *models.ConfigTLS) (credentials.TransportCredentials, error) {
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
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caPool,
	}

	return credentials.NewTLS(config), nil
}

func GetClientCN(ctx context.Context) string {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return ""
	}

	authInfo, ok := p.AuthInfo.(credentials.TLSInfo)
	if !ok {
		return ""
	}

	if len(authInfo.State.PeerCertificates) == 0 {
		return ""
	}

	return authInfo.State.PeerCertificates[0].Subject.CommonName
}
