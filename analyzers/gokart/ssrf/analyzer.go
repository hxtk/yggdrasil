package ssrf

import (
	"reflect"

	"github.com/praetorian-inc/gokart/analyzers"
	"github.com/praetorian-inc/gokart/util"
)

var Analyzer = analyzers.SSRFAnalyzer

func init() {
	Analyzer.ResultType = reflect.TypeOf([]util.Finding{})
}
