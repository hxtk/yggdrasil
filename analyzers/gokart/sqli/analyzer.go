package sqli

import (
	"reflect"

	"github.com/praetorian-inc/gokart/analyzers"
	"github.com/praetorian-inc/gokart/util"
)

var Analyzer = analyzers.SQLInjectionAnalyzer

func init() {
	Analyzer.ResultType = reflect.TypeOf([]util.Finding{})
}
