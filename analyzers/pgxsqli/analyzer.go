// Copyright 2021 Praetorian Security, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pgxsqli

import (
	"reflect"
	"strings"

	"github.com/praetorian-inc/gokart/util"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/buildssa"
)

var Analyzer = &analysis.Analyzer{
	Name:       "pgx_sql_injection",
	Doc:        "reports when SQL injection can occur",
	Run:        sqlRun,
	Requires:   []*analysis.Analyzer{buildssa.Analyzer},
	ResultType: reflect.TypeOf([]util.Finding(nil)),
}

// grab_vulnerable_sql_functions() creates map of vulnerable functions that the scanner will check
func getVulnSqlFuncs() map[string][]string {
	return map[string][]string{
		"(*github.com/jackc/pgx/v4.Conn)": {"Exec", "Query", "QueryRow", "Prepare"},
		"(*github.com/jackc/pgx/v4.Tx)":   {"Exec", "Query", "QueryRow", "Prepare"},
		"(*github.com/jackc/pgx/v5.Conn)": {"Exec", "Query", "QueryRow", "Prepare"},
		"(*github.com/jackc/pgx/v5.Tx)":   {"Exec", "Query", "QueryRow", "Prepare"},
	}
}

// sql_run runs the path traversal analyzer
func sqlRun(pass *analysis.Pass) (interface{}, error) {
	results := []util.Finding{}
	// Builds SSA model of Go code
	ssaFuncs := pass.ResultOf[buildssa.Analyzer].(*buildssa.SSA).SrcFuncs

	// Creates call graph of function calls
	cg := make(util.CallGraph)

	// Fills in call graph
	for _, fn := range ssaFuncs {
		cg.AnalyzeFunction(fn)
	}

	// Grabs vulnerable functions to scan for
	vuln_db_funcs := getVulnSqlFuncs()

	// Iterate over every specified vulnerable package
	for pkg, funcs := range vuln_db_funcs {

		// Iterate over every specified vulnerable function per package
		for _, fn := range funcs {

			// Construct full name of function
			current_function := pkg + "." + fn

			// For SQL injections we only care about the argument that holds the query string (index 1 for normal query and index 2 for Context query)
			argIndex := 2
			if strings.Contains(current_function, "Prepare") {
				argIndex = 3
			}

			// Iterate over occurrences of vulnerable function in call graph
			for _, vulnFunc := range cg[current_function] {

				// Check if argument of vulnerable function is tainted by possibly user-controlled input
				taint_analyzer := util.CreateTaintAnalyzer(pass, vulnFunc.Fn.Pos())
				if taint_analyzer.ContainsTaint(&vulnFunc.Instr.Call, &vulnFunc.Instr.Call.Args[argIndex], cg) {
					message := "Danger: possible SQL injection detected"
					targetFunc := util.GenerateTaintedCode(pass, vulnFunc.Fn, vulnFunc.Instr.Pos())
					taintSource := taint_analyzer.TaintSource
					results = append(results, util.MakeFinding(message, targetFunc, taintSource, "CWE-89: SQL Injection"))
				}
			}
		}
	}

	return results, nil
}
