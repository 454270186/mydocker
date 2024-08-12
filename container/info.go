package container

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/454270186/mydocker/constant"
)

// Container status
const (
	RUNNING = "running"
)

// Default file path
const (
	InfoLoc         = "/var/lib/mydocker/containers/"
	InfoLocalFormat = "/var/lib/mydocker/containers/%s/"
	ConfigName      = "config.json"
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
	if err := os.MkdirAll(dirPath, constant.Perm0622); err != nil {
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

func ListContainers() {
	files, err := os.ReadDir(InfoLoc)
	if err != nil {
		fmt.Printf("error while read info dir: %v\n", err)
		return
	}

	containerInfos := make([]*Info, 0, len(files))
	for _, file := range files {
		curContainer, err := getContainerInfo(file)
		if err != nil {
			fmt.Println(err)
			continue
		}

		containerInfos = append(containerInfos, curContainer)
	}

	w := tabwriter.NewWriter(os.Stdout, 12, 1, 3, ' ', 0)
	_, err = fmt.Fprintf(w, "ID\tNAME\tPID\tSTATUS\tCOMMAND\tCREATED\n")
	if err != nil {
		fmt.Printf("error while Fprintf: %v\n", err)
	}

	for _, containerInfo := range containerInfos {
		_, err = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			containerInfo.Id,
			containerInfo.Name,
			containerInfo.Pid,
			containerInfo.Status,
			containerInfo.Command,
			containerInfo.CreatedTime)
		if err != nil {
			fmt.Printf("error while Fprintf %v\n", err)
		}
	}

	if err := w.Flush(); err != nil {
		fmt.Printf("error while Flush: %v", err)
	}
}

func getContainerInfo(file os.DirEntry) (*Info, error) {
	configDirName := fmt.Sprintf(InfoLocalFormat, file.Name())
	configFileName := path.Join(configDirName, ConfigName)

	content, err := os.ReadFile(configFileName)
	if err != nil {
		return nil, fmt.Errorf("error while read config file %s: %v", configFileName, err)
	}

	info := new(Info)
	if err := json.Unmarshal(content, info); err != nil {
		return nil, fmt.Errorf("error while unmarshal config: %v", err)
	}

	return info, nil
}
