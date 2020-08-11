package release

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	tmpPath = "/tmp"
)

var (
	exeCmd    = "fastlane"
	subCmdRef = map[string]string{
		"GOOGLEPLAY": "supply",
		"APPCONNECT": "deliver",
		// "MAC":     "deliver",
	}
	initCmd = "init"
)

// NewCmd new fastlane cmd
func NewCmd(platform string, packageName string, initiate bool) (cmd *exec.Cmd, err error) {
	log.Printf("platform=%s, packageName=%s, initiate=%v", platform, packageName, initiate)

	platform = strings.ToUpper(platform)
	log.Printf("new platform=%s", platform)

	executePath, err := exec.LookPath(exeCmd)
	if err != nil {
		log.Println("fastlane not installed=", err)
	}
	fmt.Printf("fastlane is available at %s\n", executePath)

	var args []string
	args = append(args, executePath)

	subCmd, ok := subCmdRef[platform]
	if !ok {
		log.Println("platform not supported yet.")
	}
	args = append(args, subCmd)

	// path to store metadata
	s := []string{tmpPath, strings.ToLower(platform), packageName}
	metaPath := strings.Join(s, "/")
	createFilePaths(metaPath)

	if initiate {
		args = append(args, initCmd)
	}

	cmd = &exec.Cmd{
		Path: executePath,
		Args: args,
	}

	// change work dir
	cmd.Dir = metaPath

	// set cmd env
	// cmd.Env = append(os.Environ(),
	// 	"FOO=duplicate_value", // ignored
	// 	"FOO=actual_value",    // this value is used
	// )

	fmt.Printf("cmd=%s\n", cmd.String())

	return
}

func createFilePaths(path string) {
	if _, err := os.Stat(path); err == nil {
		// path/to/whatever exists
		log.Println("existed")
	} else if os.IsNotExist(err) {
		log.Println("not existed, create")
		os.MkdirAll(path, 0755)
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		log.Println("unknow error")
	}
}
