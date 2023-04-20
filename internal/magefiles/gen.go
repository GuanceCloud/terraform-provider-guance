//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Gen mg.Namespace

// Doc generate the documentation
func (ns Gen) Doc() error {
	return sh.Run("go", "generate", "./...")
}
