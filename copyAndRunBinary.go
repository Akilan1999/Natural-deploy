package natural_deploy

import (
	"github.com/melbahja/goph"
	"strings"
)

// CopyAndRunBinary Copies the binary to the tmp directory and runs the binary from there
func CopyAndRunBinary(binaryPath string, binaryName string, client *goph.Client, OS string) (string, error) {
	if strings.Contains(OS, "Windows") {
		// Probably copies to the home directory
		err := client.Upload(binaryPath, "")
		if err != nil {
			return "", err
		}
	} else if strings.Contains(OS, "Android") {
		err := client.Upload(binaryPath, "/data/local/tmp")
		if err != nil {
			return "", err
		}
	} else {
		err := client.Upload(binaryPath, "/tmp")
		if err != nil {
			return "", err
		}
	}

	// Run the binary
	run, err := client.Run("./" + binaryName)
	if err != nil {
		return "", err
	}

	// Returns output of the binary
	return string(run[:]), nil
}
