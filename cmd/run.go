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

// run is the main function that generates the stubs for the project
func run(ctx context.Context) error {
	eg, _ := errgroup.WithContext(ctx)
	eg.SetLimit(10)
	eg.Go(func() error {
		return exec.Command("stubgen", "-p", "src", "-o", "stubs").Run()
	})
	if _, err := os.Stat("stubs"); os.IsNotExist(err) {
		os.Mkdir("stubs", 0755)
	}
	fmt.Println("Generating stubs for project...")
	which, err := exec.Command("which", "pip").Output()
	if err != nil {
		return fmt.Errorf("error getting pip path: %v", err)
	}
	fmt.Printf("which pip: %s\n", strings.TrimSpace(string(which)))
	out, err := exec.Command("pip", "freeze").Output()
	if err != nil {
		return fmt.Errorf("error getting list of installed packages: %v", err)
	}
	packages := strings.Split(strings.TrimSpace(string(out)), "\n")
	for _, pkg := range packages {
		pkgCopy := pkg
		eg.Go(func() error {
			pkgName := strings.Split(pkgCopy, "==")[0]
			err := exec.Command("stubgen", "-o", "stubs/", "-m", pkgName).Run()
			if err != nil {
				log.Printf("error generating stubs for %s: %v", pkgName, err)
			} else {
				fmt.Printf("Generating stubs for %s...\n", pkgName)
			}
			err = exec.Command("mypy", "-o", "stubs/", "-m", pkgName).Run()
			if err != nil {
				return fmt.Errorf("error generating mypy stubs for %s: %v", pkgName, err)
			}
			fmt.Printf("Generating mypy stubs for %s...\n", pkgName)
			return nil
		})
	}
	return eg.Wait()
}
