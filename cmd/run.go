package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/sync/errgroup"
)

func run(ctx context.Context) error {
	// Create a new error group
	eg, _ := errgroup.WithContext(ctx)

	// Generate stubs for the project
	eg.Go(func() error {
		return exec.Command("stubgen", "-p", "src", "-o", "stubs").Run()
	})

	// Create a directory for stubs if it doesn't already exist
	eg.Go(func() error {
		if _, err := os.Stat("stubs"); os.IsNotExist(err) {
			return os.MkdirAll("stubs", 0755)
		}
		return nil
	})

	// Echo messages
	eg.Go(func() error {
		fmt.Println("Generating stubs for project...")
		out, err := exec.Command("which", "pip").Output()
		if err != nil {
			return err
		}
		fmt.Printf("which pip: %s\n", strings.TrimSpace(string(out)))
		return nil
	})

	// Get a list of installed packages using pip freeze
	out, err := exec.Command("pip", "freeze").Output()
	if err != nil {
		return err
	}
	packages := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, pkg := range packages {
		// Generate stubs for each package
		eg.Go(func() error {
			pkgName := strings.Split(pkg, "==")[0]
			err := exec.Command("stubgen", "-o", "stubs/", "-m", pkgName).Run()
			if err != nil {
				log.Printf("Error generating stubs for %s: %v", pkgName, err)
			} else {
				fmt.Printf("Generating stubs for %s...\n", pkgName)
			}
			return err
		})
	}
	fmt.Println("Stub generation complete.")
	return eg.Wait()
}
