package handler

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/454270186/mydocker/cgroups"
	"github.com/454270186/mydocker/cgroups/resource"
	"github.com/454270186/mydocker/container"
	"github.com/spf13/cobra"
)

// flags of run command
var (
	IsTTY       bool

	// cgrpup limit
	MemoryLimit string
	CpuLimit    int
	CpusetLimit string

	// volume
	Volume string

	// container process detach
	IsDetach bool

	// container name
	ContainerName string
)

func RunCmdHandler(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		log.Println("missing container command")
		return
	}

	// tty and detach cannot take effect simultaneously
	if IsTTY && IsDetach {
		log.Println("tty and detach cannot take effect simultaneously")
		return
	}

	resConf := &resource.ResourceConfig{
		MemoryLimit: MemoryLimit,
		CpuCfsQuota: CpuLimit,
		CpuSet:      CpusetLimit,
	}

	Run(IsTTY, args, resConf, Volume, ContainerName)
}

func Run(tty bool, cmdArr []string, resConf *resource.ResourceConfig, volume string, containerName string) {
	containerId := container.GetUUID()
	
	parent, writePipe := container.NewParentProcess(tty, volume)
	if parent == nil {
		return
	}
	if err := parent.Start(); err != nil {
		log.Println(err)
	}

	// record container info
	err := container.RecordContainerInfo(parent.Process.Pid, cmdArr, containerName, containerId)
	if err != nil {
		log.Println(err)
	}

	// create cgroup manager
	// - set resource limit
	// - apply container pid to cgroup
	cgroupManager := cgroups.NewCgroupManager("mydocker-cgroup")
	defer cgroupManager.Destroy()
	_ = cgroupManager.Set(resConf)
	err = cgroupManager.Apply(parent.Process.Pid)
	if err != nil {
		fmt.Println(err)
	}

	sendInitCommand(cmdArr, writePipe)

	// 如果是tty，那么父进程等待，就是前台运行，否则就是跳过，实现后台运行
	if tty {
		_ = parent.Wait()
		container.DeleteWorkSpace("/root/", volume)
		container.DeleteContainerInfo(containerId)
	}
}

// send init command to child process through pipe
func sendInitCommand(cmdArr []string, writePipe *os.File) {
	commands := strings.Join(cmdArr, " ")
	log.Printf("command all is %s\n", commands)

	writePipe.WriteString(commands)
	writePipe.Close()
}
