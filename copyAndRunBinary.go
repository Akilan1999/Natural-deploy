package natural_deploy

import (
	"github.com/melbahja/goph"
)

// CopyAndRunBinary Copies the binary to the tmp directory and runs the binary from there
func CopyAndRunBinary(binaryPath string, binaryName string, client *goph.Client, OS string) (string, error) {
	//if strings.Contains(OS, "Windows") {
	// Probably copies to the home directory
	err := client.Upload(binaryPath+binaryName, binaryName)
	if err != nil {
		return "", err
	}

	// Run the binary
	run, err := client.Run("chmod +x " + binaryName)
	if err != nil {
		return "", err
	}

	run, err = client.Run("./" + binaryName + " &")
	if err != nil {
		return "", err
	}

	// Returns output of the binary
	return string(run[:]), nil
}
