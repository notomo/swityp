package swityp_test

import (
	"testing"

	"github.com/notomo/swityp"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	tests := []struct {
		pkgPattern string
		target     string
	}{
		{
			pkgPattern: "simple",
			target:     "simple.NewType",
		},
	}

	dir := analysistest.TestData()
	for _, test := range tests {
		t.Run(test.pkgPattern, func(t *testing.T) {
			analyzer := swityp.Analyzer
			analyzer.Flags.Set("target", test.target)
			analyzer.Flags.Set("env", "GOPATH="+dir)
			analyzer.Flags.Set("env", "GO111MODULE=off")

			analysistest.Run(t, dir, analyzer, test.pkgPattern)
		})
	}
}
