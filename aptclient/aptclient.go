package aptclient

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// AptClient is cool
type AptClient struct {
	aptGetExecutable string
}

// New is cool
func New() *AptClient {
	aptGetExecutable, err := exec.LookPath("apt-get")
	if err != nil {
		log.Fatalf("could not get apt-get executable. err: %v", err)
	}

	return &AptClient{
		aptGetExecutable: aptGetExecutable,
	}
}

// Autoremove is cool
func (apt *AptClient) Autoremove() error {
	var out bytes.Buffer
	var outErr bytes.Buffer

	cmdUpgrade := &exec.Cmd{
		Path:   apt.aptGetExecutable,
		Args:   []string{apt.aptGetExecutable, "-y", "autoremove"},
		Stdout: &out,
		Stderr: &outErr,
		Env:    []string{"DEBIAN_FRONTEND=noninteractive"},
	}

	if err := cmdUpgrade.Run(); err != nil {
		return fmt.Errorf("could not run apt-get autoremove command. err: %v, out: %s, outErr: %s", err, out.String(), outErr.String())
	}

	return nil
}

// Upgrade is cool
func (apt *AptClient) Upgrade() error {
	var out bytes.Buffer
	var outErr bytes.Buffer

	cmdUpgrade := &exec.Cmd{
		Path:   apt.aptGetExecutable,
		Args:   []string{apt.aptGetExecutable, "-y", "upgrade"},
		Stdout: &out,
		Stderr: &outErr,
		Env:    []string{"DEBIAN_FRONTEND=noninteractive"},
	}

	if err := cmdUpgrade.Run(); err != nil {
		return fmt.Errorf("could not run apt-get upgrade command. err: %v, out: %s, outErr: %s", err, out.String(), outErr.String())
	}

	return nil
}

// Update is cool
func (apt *AptClient) Update() (int64, error) {
	var out bytes.Buffer
	var outErr bytes.Buffer

	cmdUpdate := &exec.Cmd{
		Path:   apt.aptGetExecutable,
		Args:   []string{apt.aptGetExecutable, "-o", "apt::cmd::show-update-stats=true", "-o", "Acquire::Check-Valid-Until=false", "-o", "Acquire::Check-Date=false", "update"},
		Stdout: &out,
		Stderr: &outErr,
		Env:    []string{"DEBIAN_FRONTEND=noninteractive"},
	}

	if err := cmdUpdate.Run(); err != nil {
		return -1, fmt.Errorf("could not run apt-get update command. err: %v, out: %s, outErr: %s", err, out.String(), outErr.String())
	}

	outString := strings.TrimSpace(out.String())
	if strings.HasSuffix(outString, "All packages are up to date.") {
		return 0, nil
	}

	updateCountRegex := regexp.MustCompile(`\n(\d+)\spackages\scan\sbe\supgraded`)
	findings := updateCountRegex.FindStringSubmatch(outString)
	if len(findings) != 2 {
		return -1, fmt.Errorf("could not find update count number in output. out: %s, regex: %s", out.String(), updateCountRegex.String())
	}

	amount, err := strconv.ParseInt(findings[1], 10, 64)
	if err != nil {
		return -1, fmt.Errorf("could not convert count from output to int. count: %s, err: %v", findings[1], err)
	}
	return amount, nil
}
