package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"time"
)

type AuthCredentials struct {
	User     string
	Password string
}

type TLSConfig struct {
	Enabled               bool
	ServerName            string
	CAFile                string
	CertFile              string
	KeyFile               string
	InsecureSkipTLSVerify bool
}

type client struct {
	client          http.Client
	baseUrl         string
	authCredentials *AuthCredentials
	tlsConfig       *TLSConfig
}

func NewClient(baseUrl string, authCredentials *AuthCredentials, tlsConfig *TLSConfig) Client {
	httpClient := http.Client{
		Timeout: 3 * time.Second,
	}

	if tlsConfig != nil && tlsConfig.Enabled {
		tlsCfg := &tls.Config{
			ServerName:         tlsConfig.ServerName,
			InsecureSkipVerify: tlsConfig.InsecureSkipTLSVerify,
		}

		if tlsConfig.CAFile != "" || tlsConfig.CertFile != "" {
			tlsCfg.ClientCAs = x509.NewCertPool()
			tlsCfg.RootCAs = x509.NewCertPool()

			if tlsConfig.CAFile != "" {
				caData, err := os.ReadFile(tlsConfig.CAFile)
				if err != nil {
					panic(fmt.Sprintf("failed to read CA file %s: %v", tlsConfig.CAFile, err))
				}
				if !tlsCfg.RootCAs.AppendCertsFromPEM(caData) {
					panic(fmt.Sprintf("failed to parse CA certificate from %s", tlsConfig.CAFile))
				}
			}

			if tlsConfig.CertFile != "" && tlsConfig.KeyFile != "" {
				cert, err := tls.LoadX509KeyPair(tlsConfig.CertFile, tlsConfig.KeyFile)
				if err != nil {
					panic(fmt.Sprintf("failed to load client cert/key: %v", err))
				}
				tlsCfg.Certificates = []tls.Certificate{cert}
			}
		}

		httpClient.Transport = &http.Transport{
			TLSClientConfig: tlsCfg,
		}
	}

	return &client{
		client:          httpClient,
		baseUrl:         baseUrl,
		authCredentials: authCredentials,
		tlsConfig:       tlsConfig,
	}
}
