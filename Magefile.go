//go:build mage
// +build mage

package main

import (
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var services = []string{"analyze", "issue"}

// Runs all process for build
func Build() error {
	mg.Deps(Protoc)
	mg.Deps(Wire)
	mg.Deps(Gqlgen)
	mg.Deps(Docker)
	return nil
}

// Runs protoc
func Protoc() error {
	os.Chdir("pb")
	defer os.Chdir("..")
	for _, t := range services {
		target := "gdean" + t + ".proto"
		if err := sh.Run(
			"protoc",
			"--go_out=.",
			"--go_opt=paths=source_relative",
			"--go-grpc_out=.",
			"--go-grpc_opt=paths=source_relative",
			target); err != nil {
			return err
		}
	}
	return nil
}

// Runs wire
func Wire() error {
	os.Chdir("cmd")
	defer os.Chdir("..")
	for _, t := range services {
		os.Chdir(t + "api")
		if err := sh.Run("wire"); err != nil {
			return err
		}
		os.Chdir("..")
	}
	return nil
}

// Runs gqlgen
func Gqlgen() error {
	os.Chdir("gql")
	defer os.Chdir("..")
	if err := sh.Run("gqlgen"); err != nil {
		return err
	}
	return nil
}

// Runs docker build
func Docker() error {
	for _, t := range append(services, "gql") {
		if err := sh.Run("docker", "build", "-t", "gdean-"+t, "-f", "Dockerfile."+t, "."); err != nil {
			return err
		}
	}
	return nil
}
