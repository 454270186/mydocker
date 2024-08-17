package container

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"syscall"

	"github.com/454270186/mydocker/constant"
)

func StopContainer(containerId string) error {
	containerInfo, err := getContainerInfoById(containerId)
	if err != nil {
		return err
	}

	pidInt, err := strconv.Atoi(containerInfo.Pid)
	if err != nil {
		return fmt.Errorf("error while convert pid type to int: %v", err)
	}

	// kill container process
	if err := syscall.Kill(pidInt, syscall.SIGTERM); err != nil {
		return fmt.Errorf("error while kill container process %d: %v", pidInt, err)
	}

	// modify container status and pid info
	containerInfo.Status = STOP
	containerInfo.Pid = ""
	newContainerInfo, err := json.Marshal(containerInfo)
	if err != nil {
		return fmt.Errorf("error while marshal new container info: %v", err)
	}

	dirPath := fmt.Sprintf(InfoLocalFormat, containerId)
	configFilePath := path.Join(dirPath, ConfigName)
	if err := os.WriteFile(configFilePath, newContainerInfo, constant.Perm0622); err != nil {
		return fmt.Errorf("error while update new container info: %v", err)
	}

	return nil
}

func getContainerInfoById(containerId string) (*Info, error) {
	dirPath := fmt.Sprintf(InfoLocalFormat, containerId)
	configFilePath := path.Join(dirPath, ConfigName)

	bytesContent, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("error while read container info file: %v", err)
	}

	var containerInfo Info
	if err := json.Unmarshal(bytesContent, &containerInfo); err != nil {
		return nil, fmt.Errorf("error while unmarshal container info: %v", err)
	}

	return &containerInfo, nil
}