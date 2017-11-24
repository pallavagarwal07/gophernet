package main

import (
	"fmt"
	"testing"
)

const PORT = "19195"

// Copied from testing.TB: The testing.TB version has a private
// method in the signature. This TB will always thus be a subset
// of testing.TB interface (tested statically below).
type TB interface {
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fail()
	FailNow()
	Failed() bool
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	Name() string
	Skip(args ...interface{})
	SkipNow()
	Skipf(format string, args ...interface{})
	Skipped() bool
	Helper()
}

// Static assertion
var _ TB = testing.TB(nil)

type customT struct {
	name    string
	failed  bool
	skipped bool
	skipStr string
	log     string
	err     string
	parent  *myT
}

var _ TB = &customT{}

func (c *customT) Stop() {
	c.parent.Stop()
}

func (*customT) Helper() {
	return
}

func (c *customT) Skipped() bool {
	return c.skipped
}

func (c *customT) Skipf(format string, args ...interface{}) {
	c.skipped = true
	c.skipStr = fmt.Sprintf(format, args...)
}

func (c *customT) SkipNow() {
	c.skipped = true
}

func (c *customT) Skip(args ...interface{}) {
	c.skipped = true
	c.skipStr = fmt.Sprint(args...)
}

func (c *customT) Name() string {
	return c.name
}

func (c *customT) Log(args ...interface{}) {
	c.log += "\n" + fmt.Sprint(args...)
}

func (c *customT) Logf(format string, args ...interface{}) {
	c.log += "\n" + fmt.Sprintf(format, args...)
}

func (c *customT) Fail() {
	c.failed = true
}

func (c *customT) FailNow() {
	c.failed = true
	c.Stop()
}

func (c *customT) Failed() bool {
	return c.failed
}

func (c *customT) Error(args ...interface{}) {
	c.failed = true
	c.err = fmt.Sprint(args...)
}

func (c *customT) Errorf(format string, args ...interface{}) {
	c.failed = true
	c.err = fmt.Sprintf(format, args...)
}

func (c *customT) Fatal(args ...interface{}) {
	c.failed = true
	c.err = fmt.Sprint(args...)
	c.Stop()
}

func (c *customT) Fatalf(format string, args ...interface{}) {
	c.failed = true
	c.err = fmt.Sprintf(format, args...)
	c.Stop()
}

type myT struct {
	TB
	backup TB
}

func (m *myT) Stop() {
	m.backup = m.TB
	m.TB = &TBnop{}
}

type TBnop struct{}

func (t *TBnop) Error(args ...interface{})                 {}
func (t *TBnop) Errorf(format string, args ...interface{}) {}
func (t *TBnop) Fail()                                     {}
func (t *TBnop) FailNow()                                  {}
func (t *TBnop) Failed() bool                              { return true }
func (t *TBnop) Fatal(args ...interface{})                 {}
func (t *TBnop) Fatalf(format string, args ...interface{}) {}
func (t *TBnop) Log(args ...interface{})                   {}
func (t *TBnop) Logf(format string, args ...interface{})   {}
func (t *TBnop) Name() string                              { return "" }
func (t *TBnop) Skip(args ...interface{})                  {}
func (t *TBnop) SkipNow()                                  {}
func (t *TBnop) Skipf(format string, args ...interface{})  {}
func (t *TBnop) Skipped() bool                             { return false }
func (t *TBnop) Helper()                                   {}
