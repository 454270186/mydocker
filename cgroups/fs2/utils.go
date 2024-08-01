package fs2

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/454270186/mydocker/constant"
)

// return the absolute path of given cgroup in file system
func getCgroupPath(cgroupPath string, autoCreated bool) (string, error) {
	cgroupRoot := defaultCgroupsFs2MountPoint
	absPath := path.Join(cgroupRoot, cgroupPath)
	if !autoCreated {
		return absPath, nil
	}

	_, err := os.Stat(absPath)
	if err != nil && os.IsNotExist(err) {
		// if cgroup path does not exist, create it
		err = os.Mkdir(absPath, constant.Perm0755)
		return absPath, err
	}

	return absPath, err
}

func applyCgroup(pid int, cgroupPath string) error {
	subCgroupPath, err := getCgroupPath(cgroupPath, true)
	if err != nil {
		return fmt.Errorf("error while get cgroup path: %v", err)
	}

	err = os.WriteFile(path.Join(subCgroupPath, "cgroup.procs"), []byte(strconv.Itoa(pid)), constant.Perm0644)
	if err != nil {
		return fmt.Errorf("error while apply process to cgroup %v: %v", cgroupPath, err)
	}

	return nil
}