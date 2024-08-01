package fs2

import "github.com/454270186/mydocker/cgroups/resource"

// the local list of all three subsystems
var LocalSubsystems = []resource.Subsystem{
	&CpuSubsystem{},
	&CpusetSubsystem{},
	&MemorySubsystem{},
}