//go:build mage
// +build mage

package main

import (
	"os"
	"path"

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

func protocWith(version string) error {
	cur, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	os.Chdir(path.Join("pb"))
	defer os.Chdir(cur)
	for _, t := range append(services, "util") {
		target := path.Join(version, "gdean"+t+".proto")
		if err := sh.Run(
			"protoc",
			"--go_out=.",
			"--go_opt=paths=source_relative",
			"--go-grpc_out=.",
			"--go-grpc_opt=paths=source_relative",
			"--experimental_allow_proto3_optional",
			target); err != nil {
			return err
		}
	}
	return nil
}

// Runs protoc
func Protoc() error {
	for _, v := range []string{"", "v1"} {
		if err := protocWith(v); err != nil {
			return nil
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

// Runs temp server via go run, support appname APP={service name} mage RunS
func RunS() error {
	api := path.Join("cmd", os.Getenv("APP")+"api")
	maingo := path.Join(api, "main.go")
	wirego := path.Join(api, "wire_gen.go")
	if err := sh.Run("go", "run", maingo, wirego); err != nil {
		return err
	}
	return nil
}
