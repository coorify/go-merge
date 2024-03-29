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
	F string
}

type source0 struct {
	A string
	B int
}

func (s source0) E() (string, error) {
	// return "", fmt.Errorf("EEROR")
	return s.A + "E", nil
}

func (s *source0) F() (string, error) {
	// return "", fmt.Errorf("EEROR")
	return s.A + "F", nil
}

type source1 struct {
	B int
	C struct {
		D *string
	}
}

func main() {
	str1 := "STR01"
	t := &target{}
	s0 := &source0{A: "A0", B: 0}
	s1 := &source1{B: 1, C: struct{ D *string }{D: &str1}}

	if err := merge.Merge(t, s0, s1); err != nil {
		panic(err)
	}

	fmt.Printf("%#v", t)
}
