//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Gen mg.Namespace

// Doc run generator over the documentation
func (ns Gen) Doc() error {
	return sh.Run(
		"go", "run",
		"github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs",
		"generate",
		"--provider-name", "guance",
		"--examples-dir", "examples",
	)
}
