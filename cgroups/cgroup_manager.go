package cgroups

import (
	"fmt"

	"github.com/454270186/mydocker/cgroups/fs2"
	"github.com/454270186/mydocker/cgroups/resource"
)

// cgroup v2 manager
type CgroupManager struct {
	Path       string
	Resource   *resource.ResourceConfig
	Subsystems []resource.Subsystem
}

func NewCgroupManager(path string) *CgroupManager {
	return &CgroupManager{
		Path: path,
		Subsystems: fs2.LocalSubsystems,
	}
}

// Apply adds given process to this cgroup
func (c *CgroupManager) Apply(pid int) error {
	fmt.Println(pid, c.Path)
	for _, SubSysInstance := range c.Subsystems {
		err := SubSysInstance.Apply(c.Path, pid)
		if err != nil {
			return fmt.Errorf("error while apply subsystem %v: %v", SubSysInstance.Name(), err)
		}
	}

	return nil
}

// Set sets resource limit for this cgroup
func (c *CgroupManager) Set(res *resource.ResourceConfig) error {
	for _, SubSysInstance := range c.Subsystems {
		err := SubSysInstance.Set(c.Path, res)
		if err != nil {
			return fmt.Errorf("error while set resouce limit for subsystem %v: %v", SubSysInstance.Name(), err)
		}
	}

	return nil
}

// Destroy releases this cgroup
func (c *CgroupManager) Destroy() error {
	for _, SubSysInstance := range c.Subsystems {
		err := SubSysInstance.Remove(c.Path)
		if err != nil {
			return fmt.Errorf("error while remove subsystem %v: %v", SubSysInstance.Name(), err)
		}
	}

	return nil
}
