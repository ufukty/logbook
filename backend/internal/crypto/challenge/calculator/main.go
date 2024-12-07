package main

import (
	"encoding/json"
	"fmt"
	"iter"
	"os"
	"slices"

	"golang.org/x/exp/maps"
)

func seq(start, stop, increment int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := start; i < stop; i += increment {
			if !yield(i) {
				return
			}
		}
	}
}

func pow(base, exponent int) int {
	if exponent == 0 {
		return 1
	}
	p := base
	for i := 2; i <= exponent; i++ {
		p *= base
	}
	return p
}

func difficulty(a, x int) int {
	return pow(a, x)
}

type params struct {
	A int `json:"a"`
	X int `json:"x"`
}

func calc() map[int]params {
	ds := map[int]params{}
	for a := range seq(2, 62, 1) {
		for x := range seq(2, 10, 1) {
			d := difficulty(a, x)
			if _, ok := ds[d]; !ok {
				ds[d] = params{A: a, X: x}
			}
		}
	}
	return ds
}

type paramsE struct {
	params
	D int `json:"d"`
}

func sortAndLinearize(ds map[int]params) []paramsE {
	ks := maps.Keys(ds)
	slices.Sort(ks)
	s := []paramsE{}
	for _, k := range ks {
		v := ds[k]
		s = append(s, paramsE{D: k, params: v})
	}
	return s
}

func print(ss []paramsE) error {
	f, err := os.Create("dataset/values.json")
	if err != nil {
		return fmt.Errorf("creating: %w", err)
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(ss)
	if err != nil {
		return fmt.Errorf("encoding: %w", err)
	}
	return nil
}

func Main() error {
	ds := calc()
	ss := sortAndLinearize(ds)
	err := print(ss)
	if err != nil {
		return fmt.Errorf("printing: %w", err)
	}
	return nil
}

func main() {
	if err := Main(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
