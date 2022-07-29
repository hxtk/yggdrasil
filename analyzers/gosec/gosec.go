package gosec

import (
	"fmt"
	"go/token"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/securego/gosec"
	"github.com/securego/gosec/rules"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/packages"

	"github.com/hxtk/yggdrasil/analyzers/pkg/result"
)

const Name = "gosec"

const ValidateTests = false // Do not check for the security of test files.

var Analyzer = &analysis.Analyzer{
	Name: Name,
	Doc:  "Inspects source code for security problems",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	gasConfig := gosec.NewConfig()
	enabledRules := rules.Generate()
	logger := log.New(ioutil.Discard, "", 0)
	analyzer := gosec.NewAnalyzer(gasConfig, ValidateTests, logger)
	analyzer.LoadRules(enabledRules.Builders())

	pkg := &packages.Package{
		Fset:      pass.Fset,
		TypesInfo: pass.TypesInfo,
		Syntax:    pass.Files,
		Types:     pass.Pkg,
	}

	analyzer.Check(pkg)
	issues, _, _ := analyzer.Report()
	if len(issues) == 0 {
		return nil, nil
	}

	for _, i := range issues {
		text := fmt.Sprintf("[%s] %s: %s", Name, i.RuleID, i.What) // TODO: use severity and confidence
		var r *result.Range
		line, err := strconv.Atoi(i.Line)
		if err != nil {
			r = &result.Range{}
			if n, rerr := fmt.Sscanf(i.Line, "%d-%d", &r.From, &r.To); rerr != nil || n != 2 {
				//lintCtx.Log.Warnf("Can't convert gosec line number %q of %v to int: %s", i.Line, i, err)
				continue
			}
			line = r.From
		}

		pass.Reportf(token.Pos(line), text)
	}

	return nil, nil
}
