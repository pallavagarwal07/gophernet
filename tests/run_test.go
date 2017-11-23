package main

import (
	"testing"
)

func TestAll(t *testing.T) {
	for name, test := range tests {
		t.Run(name, func(t *testing.T) { test(t) })
	}
}
