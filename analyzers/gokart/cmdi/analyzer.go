package command_injection

import (
	"reflect"

	"github.com/praetorian-inc/gokart/analyzers"
	"github.com/praetorian-inc/gokart/util"
)

var Analyzer = analyzers.CommandInjectionAnalyzer

func init() {
	Analyzer.ResultType = reflect.TypeOf([]util.Finding{})
}
