package main

import "sort"

type Pair struct {
	fst string
	snd string
}

func less(a Pair, b Pair) bool {
	if a.fst < b.fst {
		return true
	}
	if a.fst > b.fst {
		return false
	}
	return a.snd < b.snd
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
