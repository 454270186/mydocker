package fs2

import (
	"fmt"
	"os"
	"path"

	"github.com/454270186/mydocker/cgroups/resource"
	"github.com/454270186/mydocker/constant"
)

// cpu set subsystem
// implement System interface
type CpusetSubsystem struct{}

func (c *CpusetSubsystem) Name() string {
	return "cpuset"
}

func (c *CpusetSubsystem) Set(cgroupPath string, res *resource.ResourceConfig) error {
	if res.CpuSet == "" {
		return nil
	}

	subCgroupPath, err := getCgroupPath(cgroupPath, true)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(subCgroupPath, "cpuset.cpus"), []byte(res.CpuSet), constant.Perm0644)
	if err != nil {
		return fmt.Errorf("error while set cpuset for cgroup %v: %v", cgroupPath, err)
	}

	return nil
}

func (c *CpusetSubsystem) Apply(cgroupPath string, pid int) error {
	return applyCgroup(pid, cgroupPath)
}

func (c *CpusetSubsystem) Remove(cgroupPath string) error {
	subCgroupPath, err := getCgroupPath(cgroupPath, true)
	if err != nil {
		return err
	}

	return os.RemoveAll(subCgroupPath)
}