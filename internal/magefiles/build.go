//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Build mg.Namespace

// Install run installation on the provider to local
func (ns Build) Install() error {
	return sh.Run("go", "install", ".")
}
