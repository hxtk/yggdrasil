package server

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/hxtk/yggdrasil/common/config/tlsconfig"
)

func FromViper(v *viper.Viper) (*Server, error) {
	_, err := tlsconfig.FromViper(v.Sub("tls"))
	if err != nil {
		log.WithError(err).Fatal("Error reading TLS Config")
	}

	s := New()

	return s, nil
}
