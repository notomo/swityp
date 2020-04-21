package main

import "fmt"

type NewType string

var (
	NewTypeA = NewType("a")
	NewTypeB = NewType("b")
	NewTypeC = NewType("c")
)

type OtherNewType string

var (
	OtherNewTypeA = OtherNewType("a")
	OtherNewTypeB = OtherNewType("b")
)

func main() {
	switch typ := NewTypeA; typ { // want "non-exhaustive switch: `NewTypeC` not covered"
	case NewTypeA:
	case NewTypeB:
	}

	switch typ := OtherNewTypeA; typ {
	case OtherNewTypeA:
	}

	typ := NewTypeA
	switch typ {
	case NewTypeA:
	default:
	}

	// avoid unsing
	fmt.Println(NewTypeC)
	fmt.Println(OtherNewTypeB)
}
