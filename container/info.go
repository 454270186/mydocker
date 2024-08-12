package container

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/454270186/mydocker/constant"
)

// Container status
const (
	RUNNING = "running"
)

// Default file path
const (
	InfoLocalFormat  = "/var/run/%s/"
	ConfigName = "config.json"
)

type Info struct {
	Pid         string `json:"pid"`
	Id          string `json:"id"`
	Name        string `json:"name"`
	Command     string `json:"command"`
	CreatedTime string `json:"createdTime"`
	Status      string `json:"status"`
}

func RecordContainerInfo(containerPid int, commandArr []string, containerName, containerId string) error {
	// if not given container name, use container id as default container name
	if containerName == "" {
		containerName = containerId
	}

	command := strings.Join(commandArr, " ")
	containerInfo := Info{
		Pid:         strconv.Itoa(containerPid),
		Id:          containerId,
		Name:        containerName,
		Command:     command,
		CreatedTime: time.Now().Format("2006-01-02 15:04:05"),
		Status:      RUNNING,
	}

	infoJsonBytes, err := json.Marshal(&containerInfo)
	if err != nil {
		return fmt.Errorf("error while marshal container info: %v", err)
	}
	infoJsonStr := string(infoJsonBytes)
	
	dirPath := fmt.Sprintf(InfoLocalFormat, containerId)
	if err := os.Mkdir(dirPath, constant.Perm0622); err != nil {
		return fmt.Errorf("error while create container info dir: %v", err)
	}

	// write container file
	fileName := path.Join(dirPath, ConfigName)
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error while create container info file: %v", err)
	}
	defer file.Close()

	if _, err := file.WriteString(infoJsonStr); err != nil {
		return fmt.Errorf("error while write container info: %v", err)
	}

	return nil
}

func DeleteContainerInfo(containerId string) {
	dirPath := fmt.Sprintf(InfoLocalFormat, containerId)
	if err := os.RemoveAll(dirPath); err != nil {
		fmt.Println(err)
	}
}