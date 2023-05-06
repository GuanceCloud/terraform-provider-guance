//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"os"
)

type Test mg.Namespace

// Acc run acceptance test for specified resource
func (ns Test) Acc() error {
	envKeys := []string{"GUANCE_ACCESS_TOKEN", "GUANCE_REGION"}
	envVars := map[string]string{}
	for _, k := range envKeys {
		envVars[k] = os.Getenv(k)
	}
	os.Chdir("examples")
	defer os.Chdir("..")
	return sh.RunWithV(envVars, "go", "test", "-v", "./...")
}
