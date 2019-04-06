package helpers

import (
	// System
	"os/exec"

	// 3rd Party
	log "github.com/sirupsen/logrus"
)

type BashCmd struct {
	Cmd      string
	Args     []string
	ExecPath string
}

// Executes CLI commands
func ExecBashCmd(c *BashCmd) {

	// Set command and argument options
	cmd := c.Cmd
	cmdArgs := c.Args

	// Execute command
	cmdRun := exec.Command(cmd, cmdArgs...)
	cmdRun.Dir = c.ExecPath

	// Print stdout & stderr
	out, err := cmdRun.CombinedOutput()
	if err != nil {
		log.Errorf("\t %v", string(out))
		log.Fatalf("\t %v", err)
	}

	// Prompt

}
