package container

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

const (
	EnvExecPid = "mydocker_pid"
	EnvExecCmd = "mydocker_cmd"
)

func ExecContainer(containerId string, cmdArr []string) {
	pid, err := getPidByContainerId(containerId)
	if err != nil {
		fmt.Println(err)
		return
	}

	cmd := exec.Command("/proc/self/exe", "exec")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmdStr := strings.Join(cmdArr, " ")
	log.Printf("container pid %s, commands %s\n", pid, cmdStr)
	_ = os.Setenv(EnvExecPid, pid)
	_ = os.Setenv(EnvExecCmd, cmdStr)

	if err := cmd.Run(); err != nil {
		fmt.Printf("error while exec container %s: %v\n", containerId, err)
	}
}

func getPidByContainerId(containerId string) (string, error) {
	dirPath := fmt.Sprintf(InfoLocalFormat, containerId)
	configFilePath := path.Join(dirPath, ConfigName)

	bytesContent, err := os.ReadFile(configFilePath)
	if err != nil {
		return "", fmt.Errorf("error while read container info file [%s]: %v", configFilePath, err)
	}

	containerInfo := Info{}
	if err := json.Unmarshal(bytesContent, &containerInfo); err != nil {
		return "", fmt.Errorf("error while unmarshal container info: %v", err)
	}

	return containerInfo.Pid, nil
}
