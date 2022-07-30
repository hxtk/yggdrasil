package rsa

import (
	"reflect"

	"github.com/praetorian-inc/gokart/analyzers"
	"github.com/praetorian-inc/gokart/util"
)

var Analyzer = analyzers.RsaKeylenAnalyzer

func init() {
	Analyzer.ResultType = reflect.TypeOf([]util.Finding{})
}
