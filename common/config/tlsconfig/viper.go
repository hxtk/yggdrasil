package tlsconfig

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"

	"github.com/spf13/viper"
)

func FromViper(v *viper.Viper) (*tls.Config, error) {
	if v.IsSet("tls.enabled") && !viper.GetBool("tls.enabled") {
		return nil, nil
	}

	config := &tls.Config{}

	certFile := v.GetString("tls.certificate")
	keyFile := v.GetString("tls.key")

	certificate, err := tls.LoadX509KeyPair(
		certFile,
		keyFile,
	)
	if err != nil {
		return nil, err
	}
	config.Certificates = []tls.Certificate{certificate}

	if v.IsSet("tls.insecure") {
		config.InsecureSkipVerify = v.GetBool("tls.insecure")
	}

	if v.IsSet("tls.ca") {
		certBytes, err := ioutil.ReadFile(v.GetString("tls.ca"))
		if err != nil {
			return nil, err
		}

		certPool := x509.NewCertPool()
		if ok := certPool.AppendCertsFromPEM(certBytes); !ok {
			return nil, errors.New("error parsing CA")
		}

		config.ClientCAs = certPool
		config.RootCAs = certPool
	}

	return config, nil
}
