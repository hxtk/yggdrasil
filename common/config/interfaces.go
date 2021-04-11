package config

import (
	"github.com/spf13/pflag"
)

// FlagBinder is an interface for an object that provides configuration
// options which can be expressed as command line flags.
type FlagBinder interface {
	// AddFlags adds a FlagBinder's configuration options to `flags`.
	AddFlags(flags *pflag.FlagSet)
}
