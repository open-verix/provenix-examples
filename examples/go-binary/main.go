package main

import (
	"fmt"
	"os"
)

var (
	// Version is set via ldflags during build
	Version = "dev"
	// Commit is set via ldflags during build
	Commit = "unknown"
	// BuildDate is set via ldflags during build
	BuildDate = "unknown"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Printf("go-binary-example version %s\n", Version)
		fmt.Printf("Commit: %s\n", Commit)
		fmt.Printf("Built: %s\n", BuildDate)
		return
	}

	fmt.Println("Hello from Provenix Go Binary Example!")
	fmt.Printf("This is version %s\n", Version)
	fmt.Println("\nThis example demonstrates:")
	fmt.Println("  - Binary SBOM generation")
	fmt.Println("  - Keyless signing with GitHub OIDC")
	fmt.Println("  - Vulnerability scanning")
	fmt.Println("  - Attestation workflow")
	fmt.Println("\nRun with 'version' argument to see build info")
}
