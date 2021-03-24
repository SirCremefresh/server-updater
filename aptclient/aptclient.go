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

type UpgradeInfo struct {
	Upgraded         int
	NewInstalled     int
	ToRemove         int
	NotUpgraded      int
	UpgradedPackages []UpgradedPackage
}

type UpgradedPackage struct {
	Name        string
	FromVersion string
	ToVersion   string
}

func (upgradeInfo *UpgradeInfo) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("upgraded: %d, new installed: %d, to remove: %d, not upgraded: %d.\n",
		upgradeInfo.Upgraded, upgradeInfo.NewInstalled, upgradeInfo.ToRemove, upgradeInfo.NotUpgraded))

	if len(upgradeInfo.UpgradedPackages) > 0 {
		sb.WriteString("upgraded packages:\n")
	}

	for _, upgradedPackages := range upgradeInfo.UpgradedPackages {
		sb.WriteString(fmt.Sprintf("%s %s -> %s\n", upgradedPackages.Name, upgradedPackages.FromVersion, upgradedPackages.ToVersion))
	}

	return sb.String()
}

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
func (apt *AptClient) Upgrade() (*UpgradeInfo, error) {
	var out bytes.Buffer
	var outErr bytes.Buffer

	cmdUpgrade := &exec.Cmd{
		Path:   apt.aptGetExecutable,
		Args:   []string{apt.aptGetExecutable, "--yes", "--with-new-pkgs", "--verbose-versions", "upgrade"},
		Stdout: &out,
		Stderr: &outErr,
		Env:    []string{"DEBIAN_FRONTEND=noninteractive"},
	}

	if err := cmdUpgrade.Run(); err != nil {
		return nil, fmt.Errorf("could not run apt-get upgrade command. err: %v, out: %s, outErr: %s", err, out.String(), outErr.String())
	}

	return parseUpgradeOut(out.String())
}

func parseUpgradeOut(upgradeOut string) (*UpgradeInfo, error) {
	var re = regexp.MustCompile(`(?m)^\s+(.+)\s\((.*)\s=>\s(.+)\)`)
	var upgradedPackages = make([]UpgradedPackage, 0)
	for _, match := range re.FindAllSubmatch([]byte(upgradeOut), -1) {
		upgradedPackages = append(upgradedPackages, UpgradedPackage{
			Name:        string(match[1]),
			FromVersion: string(match[2]),
			ToVersion:   string(match[3]),
		})
	}

	upgradeCountsRegex := regexp.MustCompile(`(\d+)\supgraded,\s(\d+)\snewly\sinstalled,\s(\d+)\sto\sremove\sand\s(\d+)\snot\supgraded`)
	findings := upgradeCountsRegex.FindStringSubmatch(upgradeOut)
	if len(findings) != 5 {
		return nil, fmt.Errorf("could not find upgraded counts in output. out: %s, regex: %s", upgradeOut, upgradeCountsRegex.String())
	}

	upgradeInfo := &UpgradeInfo{
		UpgradedPackages: upgradedPackages,
	}

	upgraded, err := strconv.Atoi(findings[1])
	upgradeInfo.Upgraded = upgraded
	if err == nil {
		var newInstalled int
		newInstalled, err = strconv.Atoi(findings[2])
		upgradeInfo.NewInstalled = newInstalled
	}
	if err == nil {
		var toRemove int
		toRemove, err = strconv.Atoi(findings[3])
		upgradeInfo.ToRemove = toRemove
	}
	if err == nil {
		var notUpgraded int
		notUpgraded, err = strconv.Atoi(findings[4])
		upgradeInfo.NotUpgraded = notUpgraded
	}
	if err != nil {
		return nil, fmt.Errorf("could not parse numbers in upgraded counts in output. out: %s, regex: %s", upgradeOut, upgradeCountsRegex.String())
	}

	return upgradeInfo, nil
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
