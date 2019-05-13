package ssl

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"github.com/project-flogo/core/support/log"
)

const ConfigSchema = `
{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "type": "object",
  "properties": {
    "caFile": {
      "type": "string"
    },
    "certFile": {
      "type": "string"
    },
    "keyFile": {
      "type": "string"
    },
    "skipVerify": {
      "type": "boolean"
    },
    "useSystemCert": {
      "type": "boolean"
    }
  }
}`

type Config struct {
	CAFile        string `json:"caFile"`
	CertFile      string `json:"certFile"`
	KeyFile       string `json:"keyFile"`
	SkipVerify    bool   `json:"skipVerify"`
	UseSystemCert bool   `json:"useSystemCert"`
}

func NewClientTLSConfig(config *Config) (*tls.Config, error) {
	tlsConfig := &tls.Config{
		//MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: config.SkipVerify,
	}

	var caCertPool *x509.CertPool

	if config.UseSystemCert {
		caCertPool, _ = x509.SystemCertPool()
		if caCertPool == nil {
			log.RootLogger().Warnf("unable to get system cert pool, using empty pool")
		}
	}

	if caCertPool == nil {
		caCertPool = x509.NewCertPool()
	}

	if config.CAFile != "" {
		caCert, err := ioutil.ReadFile(config.CAFile)
		if err != nil {
			return nil, err
		}
		caCertPool.AppendCertsFromPEM(caCert)
	}

	tlsConfig.RootCAs = caCertPool

	if config.CertFile != "" && config.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(config.CertFile, config.CertFile)
		if err != nil {
			return nil, err
		}

		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	return tlsConfig, nil
}
