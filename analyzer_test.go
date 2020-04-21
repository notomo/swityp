package swityp_test

import (
	"testing"

	"github.com/notomo/swityp"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	tests := []struct {
		name   string
		target string
	}{
		{
			name:   "simple",
			target: "simple.NewType",
		},
	}

	data := analysistest.TestData()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			analyzer := swityp.Analyzer
			analyzer.Flags.Set("target", test.target)
			analysistest.Run(t, data, analyzer, test.name)
		})
	}
}
