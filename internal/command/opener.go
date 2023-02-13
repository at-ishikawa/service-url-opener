package command

import (
	"fmt"
	"os/exec"
)

func openBrowser(url string) error {
	// TODO: support multiple commands
	command := "open"

	if err := exec.Command(command, url).Run(); err != nil {
		return fmt.Errorf("exec.Command > %w", err)
	}
	return nil
}
