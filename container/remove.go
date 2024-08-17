package container

import (
	"fmt"
	"os"
)

func RemoveContainer(containerId string, isForce bool) error {
	containerInfo, err := getContainerInfoById(containerId)
	if err != nil {
		return fmt.Errorf("error while get container info by id %s: %v", containerId, err)
	}

	switch containerInfo.Status {
	case RUNNING:
		// cannot remove running container without force flag
		// if force flag is true, first stop the container then remove it
		if !isForce {
			return fmt.Errorf("cannot remove running container[%s], stop it first", containerId)
		}
		StopContainer(containerId)
		RemoveContainer(containerId, isForce)
	case STOP:
		dirPath := fmt.Sprintf(InfoLocalFormat, containerId)
		if err := os.RemoveAll(dirPath); err != nil {
			return fmt.Errorf("error while remove container %s dir: %v", containerId, err)
		}
	default:
		return fmt.Errorf("invalid container status[%s]", containerInfo.Status)
	}

	return nil
}