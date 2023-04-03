//go:build mage
// +build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Dev mg.Namespace

// Fmt format the code
func (ns Dev) All() error {
	mg.Deps(ns.Fmt, ns.Lint)
	return nil
}

// Fmt format the code
func (ns Dev) Fmt() error {
	err := sh.Run("goimports", "-w", ".")
	if err != nil {
		return fmt.Errorf("goimports failed: %w", err)
	}

	files, err := ListFileByExt(".", "go")
	if err != nil {
		return err
	}
	err = sh.Run("gofmt", append([]string{"-w", "-s"}, files...)...)
	if err != nil {
		return fmt.Errorf("format code failed: %w", err)
	}
	return nil
}

// Lint lint the code
func (ns Dev) Lint() error {
	mg.Deps(ns.Fmt)
	return sh.Run("golangci-lint", "run", "--fix", "--allow-parallel-runners")
}

type Build mg.Namespace

// Install install the provider to local
func (ns Build) Install() error {
	return sh.Run("go", "install", ".")
}

type Gen mg.Namespace

// Doc generate the documentation
func (ns Gen) Doc() error {
	return sh.Run("go", "generate", "./...")
}
