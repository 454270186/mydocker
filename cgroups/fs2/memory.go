package fs2

import (
	"fmt"
	"os"
	"path"

	"github.com/454270186/mydocker/cgroups/resource"
	"github.com/454270186/mydocker/constant"
)

// memory subsystem
// implement Subsystem interface
type MemorySubsystem struct{}

// Name returns the name of subsystem
func (m *MemorySubsystem) Name() string {
	return "memory"
}

// Set sets the momory limit for givn cgroup
func (m *MemorySubsystem) Set(cgroupPath string, res *resource.ResourceConfig) error {
	if res.MemoryLimit == "" {
		return nil
	}

	subCgroupPath, err := getCgroupPath(cgroupPath, true)
	if err != nil {
		return err
	}

	// set memory limit for thie cgroup
	err = os.WriteFile(path.Join(subCgroupPath, "memory.max"), []byte(res.MemoryLimit), constant.Perm0644)
	if err != nil {
		return fmt.Errorf("error while set memory limit for cgroup %v: %v", cgroupPath, err)
	}
	
	return nil
}

// Apply add pid to given cgroup
func (m *MemorySubsystem) Apply(cgroupPath string, pid int) error {
	return applyCgroup(pid, cgroupPath)
}

func (m *MemorySubsystem) Remove(cgroupPath string) error {
	subCgroupPath, err := getCgroupPath(cgroupPath, true)
	if err != nil {
		return err
	}

	return os.RemoveAll(subCgroupPath)
}