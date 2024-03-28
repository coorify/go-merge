package main

import (
	"fmt"

	"github.com/coorify/go-merge"
)

type target struct {
	A string
	B int
	C struct {
		D string
	}
	E string
	F int
}

type source0 struct {
	A string
	B int
}

func (s source0) E() (string, error) {
	// return "", fmt.Errorf("EEROR")
	return s.A + "E", nil
}

func (s source0) F() (string, error) {
	// return "", fmt.Errorf("EEROR")
	return s.A + "F", nil
}

type source1 struct {
	B int
	C struct {
		D string
	}
}

func main() {
	t := &target{}
	s0 := &source0{A: "A0", B: 0}
	s1 := &source1{B: 1, C: struct{ D string }{D: "D0"}}

	if err := merge.Merge(t, s0, s1); err != nil {
		panic(err)
	}

	fmt.Printf("%#v", t)
}
