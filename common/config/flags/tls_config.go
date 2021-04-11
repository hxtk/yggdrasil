package flags

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ThalesIgnite/crypto11"
	"github.com/spf13/pflag"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/hxtk/yggdrasil/common/config"
	"github.com/hxtk/yggdrasil/common/config/tlsconfig"
)

var (
	ErrInvalidConfig = fmt.Errorf("Invalid set of configuration options.")
	ErrParsingCA     = fmt.Errorf("Certificate Authority could not be parsed.")
)

type TLSRole int8

const (
	undefined_role TLSRole = iota
	SERVER_ROLE
	CLIENT_ROLE
)

type ConnectionFlags struct {
	// CertFile is a path to a PEM-Encoded X509 Certificate File.
	// This file will be used if PKCS11Provider is not provided or is
	// unable to be used, and KeyFile is also provided.
	CertFile *string

	// KeyFile is a path to a PEM-Encoded private key corresponding to
	// the leaf certificate stored in CertFile.
	KeyFile *string

	// Password is used to decrypt PEM-Encoded private keys, or as the
	// pin for PKCS11 providers.
	Password *string

	// CAFile is the path to the PEM-Encoded X509 Certificate Authority
	// that should be used to authenticate the other side of the TLS
	// connection.
	CAFile *string

	// TLSInsecure indicates whether the remote should be verified.
	TLSInsecure *bool

	// PlainText disables TLS.
	PlainText *bool

	// PKCS11Provider is a path to a PKCS11 shared library which may be
	// used to load a certificate from a hardware security module.
	// If present, this will always take precedence over Certfile and
	// KeyFile.
	PKCS11Provider *string

	// Port is the port on which the server is listening.
	Port *int

	// HostName is the DNS name or IP address at which the server may be
	// found.
	HostName *string

	// TLSHostName is the name that will be used to validate the server's
	// certificate, if it is different from the HostName.
	TLSHostName *string

	// MutualTLSRequred indicates whether the client should authenticate
	// with their own TLS certificate. In clients, this unused.
	MutualTLSRequired *bool
}

// Type assert ConnectionFlags implements tlsconfig.Provider
var _ tlsconfig.Provider = &ConnectionFlags{}

// Type assert ConnectionFlags implements config.FlagFileBinder
var _ config.FlagBinder = &ConnectionFlags{}

// NewConnectionFlags returns a default-valued ConnectionFlags object for the
// given role.
//
// In the SERVER role, the MutualTLSRequired flag will be enabled.
func NewConnectionFlags(role TLSRole) *ConnectionFlags {
	insecure := false
	plaintext := false
	port := 42913

	tmp := false
	var mutualTLSRequired *bool = nil
	if role == SERVER_ROLE {
		mutualTLSRequired = &tmp
	}

	return &ConnectionFlags{
		CertFile:          stringptr(""),
		KeyFile:           stringptr(""),
		CAFile:            stringptr(""),
		Password:          stringptr(""),
		TLSInsecure:       &insecure,
		PlainText:         &plaintext,
		PKCS11Provider:    stringptr(""),
		Port:              &port,
		HostName:          stringptr(""),
		TLSHostName:       stringptr(""),
		MutualTLSRequired: mutualTLSRequired,
	}
}

// String implements Stringer for ConnectionFlags.
func (c *ConnectionFlags) String() string {
	res := ""
	if c.CertFile != nil {
		res += fmt.Sprintf("Cert File: %s\n", *c.CertFile)
	}
	if c.KeyFile != nil {
		res += fmt.Sprintf("Key File: %s\n", *c.KeyFile)
	}
	if c.CAFile != nil {
		res += fmt.Sprintf("CA File: %s\n", *c.CAFile)
	}
	if c.Password != nil {
		res += fmt.Sprintf("Password: %s\n", *c.Password)
	}
	if c.TLSInsecure != nil {
		res += fmt.Sprintf("TLSInsecure: %t\n", *c.TLSInsecure)
	}
	if c.PlainText != nil {
		res += fmt.Sprintf("PlainText: %t\n", *c.PlainText)
	}
	if c.PKCS11Provider != nil {
		res += fmt.Sprintf("PKCS11Provider: %s\n", *c.PKCS11Provider)
	}
	if c.Port != nil {
		res += fmt.Sprintf("Port: %d\n", *c.Port)
	}
	if c.HostName != nil {
		res += fmt.Sprintf("HostName: %s\n", *c.HostName)
	}
	if c.TLSHostName != nil {
		res += fmt.Sprintf("TLSHostName: %s\n", *c.TLSHostName)
	}
	if c.MutualTLSRequired != nil {
		res += fmt.Sprintf("MutualTLSRequired: %t\n", *c.MutualTLSRequired)
	}

	return res
}

func (c *ConnectionFlags) ToFlagSet() *pflag.FlagSet {
	flags := pflag.NewFlagSet("connection-flags", pflag.ExitOnError)
	if c.CertFile != nil {
		flags.StringVarP(c.CertFile, "cert", "E", *c.CertFile, "Path of the x509 Certificate to use")
	}
	if c.KeyFile != nil {
		flags.StringVar(c.KeyFile, "key", *c.KeyFile, "Path of the private key to use")
	}
	if c.CAFile != nil {
		flags.StringVar(c.CAFile, "cacert", *c.CAFile, "Path of the trust store to use (default: system)")
	}
	if c.Password != nil {
		flags.StringVar(c.CAFile, "password", *c.Password, "Password for key-file or PKCS#11 token. If not provided, a prompt will be given.")
	}
	if c.TLSInsecure != nil {
		flags.BoolVarP(c.TLSInsecure, "tls-insecure", "k", *c.TLSInsecure, "Disable authentication of remote connections")
	}
	if c.PlainText != nil {
		flags.BoolVar(c.PlainText, "plaintext", *c.PlainText, "Disable encryption of remote connections. Implies --tls-insecure.")
	}
	if c.PKCS11Provider != nil {
		flags.StringVarP(c.PKCS11Provider, "pkcs11", "s", *c.PKCS11Provider, "PKCS11 shared library path for HSM certificate (takes precedence over cert-file)")
	}
	if c.Port != nil {
		flags.IntVarP(c.Port, "port", "p", *c.Port, "Server port number")
	}
	if c.HostName != nil {
		flags.StringVarP(c.HostName, "host", "H", *c.HostName, "Server DNS name or IP address")
	}
	if c.TLSHostName != nil {
		flags.StringVar(c.TLSHostName, "tls-hostname", *c.TLSHostName, "Server name for TLS verification, if different from host")
	}
	if c.MutualTLSRequired != nil {
		flags.BoolVar(c.MutualTLSRequired, "enforce-client-tls", *c.MutualTLSRequired, "Require clients to mutually authenticate TLS connection (implied by ca-file)")
	}

	return flags
}

// AddFlags binds connection configuration flags to a given flagset.
func (c *ConnectionFlags) AddFlags(flags *pflag.FlagSet) {
	flags.AddFlagSet(c.ToFlagSet())
}

func (c *ConnectionFlags) ToTLSConfig() (*tls.Config, error) {
	config := &tls.Config{}

	if *c.PKCS11Provider != "" {
		if *c.Password == "" {
			fmt.Printf("PKCS#11 Pin: ")
			tmp, err := terminal.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return nil, err
			}
			*c.Password = string(tmp)
		}

		ctx, err := crypto11.Configure(&crypto11.Config{
			Path: *c.PKCS11Provider,
			Pin:  *c.Password,
		})
		if err != nil {
			return nil, err
		}

		certs, err := ctx.FindAllPairedCertificates()
		if err != nil {
			return nil, err
		}

		config.Certificates = certs
	} else if *c.CertFile != "" && *c.KeyFile != "" {
		certificate, err := tls.LoadX509KeyPair(*c.CertFile, *c.KeyFile)
		if err != nil {
			return nil, err
		}

		config.Certificates = []tls.Certificate{certificate}
	} else {
		return nil, ErrInvalidConfig
	}

	config.InsecureSkipVerify = *c.TLSInsecure
	if c.MutualTLSRequired != nil && *c.MutualTLSRequired {
		config.ClientAuth = tls.RequireAndVerifyClientCert
		if *c.TLSInsecure {
			config.ClientAuth = tls.RequireAnyClientCert
		}
	}

	if *c.CAFile != "" {
		certBytes, err := ioutil.ReadFile(*c.CAFile)
		if err != nil {
			return nil, err
		}

		certPool := x509.NewCertPool()
		if ok := certPool.AppendCertsFromPEM(certBytes); !ok {
			return nil, ErrParsingCA
		}

		config.ClientCAs = certPool
		config.RootCAs = certPool
	}

	return config, nil
}

func stringptr(val string) *string {
	return &val
}
