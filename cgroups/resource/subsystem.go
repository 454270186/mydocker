package resource

type ResourceConfig struct {
	MemoryLimit string
	CpuCfsQuota  int
	CpuSet       string
}

type Subsystem interface {
	// return subsystem name
	Name() string
	// Set given resource config for the cgroup which is determined by path
	Set(path string, res *ResourceConfig) error
	// Apply process to the cgroup which is determined by path
	Apply(path string, pid int) error
	// Remove the cgroup determined by path
	Remove(path string) error
}
