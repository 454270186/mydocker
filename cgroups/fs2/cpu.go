package fs2

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/454270186/mydocker/cgroups/resource"
	"github.com/454270186/mydocker/constant"
)

const (
	DefaultPeriod = 100000
)

// cpu subsystem
// implement Subsystem interface
type CpuSubsystem struct{}

// Name returns the name of subsystem
func (c *CpuSubsystem) Name() string {
	return "cpu"
}

// Set sets the momory limit for givn cgroup
func (c *CpuSubsystem) Set(cgroupPath string, res *resource.ResourceConfig) error {
	if res.CpuCfsQuota == 0 {
		return nil
	}

	subCgroupPath, err := getCgroupPath(cgroupPath, true)
	if err != nil {
		return err
	}

	// set cpu使用率
	err = os.WriteFile(path.Join(subCgroupPath, "cpu.max"), []byte(fmt.Sprintf("%v %v", strconv.Itoa(DefaultPeriod/100*res.CpuCfsQuota), DefaultPeriod)), constant.Perm0644)
	if err != nil {
		return fmt.Errorf("error while set cpu share for cgroup %v: %v", cgroupPath, err)
	}
	
	return nil
}

// Apply add pid to given cgroup
func (c *CpuSubsystem) Apply(cgroupPath string, pid int) error {
	return applyCgroup(pid, cgroupPath)
}

func (c *CpuSubsystem) Remove(cgroupPath string) error {
	subCgroupPath, err := getCgroupPath(cgroupPath, true)
	if err != nil {
		return err
	}

	return os.RemoveAll(subCgroupPath)
}