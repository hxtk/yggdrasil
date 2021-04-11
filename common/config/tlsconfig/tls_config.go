package tlsconfig

import (
	"crypto/tls"
)

type Provider interface {
	ToTLSConfig() (*tls.Config, error)
}
