package main

import (
	"github.com/notomo/swityp"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(swityp.Analyzer)
}
