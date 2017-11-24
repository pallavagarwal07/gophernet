package main

import (
	"math/rand"
	"sort"
)

type Pair [2]string

func less(a Pair, b Pair) bool {
	if a[0] < b[0] {
		return true
	}
	if a[0] > b[0] {
		return false
	}
	return a[1] < b[1]
}

func sortPair(k []Pair) []Pair {
	sort.Slice(k, func(i, j int) bool {
		return less(k[i], k[j])
	})
	return k
}

type Request struct {
	Method string // GET or POST
	Params []Pair
}

func deterministicBinData() []byte {
	rand.Seed(10) // Set to a deterministic value everytime.
	arr := make([]byte, 100, 100)
	n, err := rand.Read(arr)
	if err != nil {
		panic(err)
	}
	return arr[:n]
}
