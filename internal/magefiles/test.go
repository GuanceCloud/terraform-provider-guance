//go:build mage
// +build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Test mg.Namespace

// AccOne run acceptance test for specified resource
func (ns Test) AccOne(name string) error {
	return sh.RunWithV(
		map[string]string{
			"TF_ACC": "1",
			"TF_LOG": "INFO",
		},
		"go", "test", "-v", "-count", "1", "-run", fmt.Sprintf("^TestAcc%s$", name), "./tests/...",
	)
}
