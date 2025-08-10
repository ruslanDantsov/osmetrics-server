package main

import (
	"github.com/ultraware/funlen"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"

	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/unusedresult"

	"honnef.co/go/tools/staticcheck"

	"github.com/kisielk/errcheck/errcheck"
)

func main() {
	var analyzers []*analysis.Analyzer

	analyzers = append(analyzers,
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		unusedresult.Analyzer,
	)

	for _, a := range staticcheck.Analyzers {
		if a.Analyzer.Name[:2] == "SA" || a.Analyzer.Name == "ST1000" {
			analyzers = append(analyzers, a.Analyzer)
		}
	}

	analyzers = append(analyzers, errcheck.Analyzer)
	analyzers = append(analyzers, funlen.NewAnalyzer(100, 50, true))

	analyzers = append(analyzers, Analyzer)

	multichecker.Main(analyzers...)
}
