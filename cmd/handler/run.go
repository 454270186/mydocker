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

var (
	IsTTY       bool
	MemoryLimit string
	CpuLimit    int
	CpusetLimit string
)

func RunCmdHandler(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		log.Println("missing container command")
		return
	}

	resConf := &resource.ResourceConfig{
		MemoryLimit: MemoryLimit,
		CpuCfsQuota: CpuLimit,
		CpuSet:      CpusetLimit,
	}

	Run(IsTTY, args, resConf)
}

func Run(tty bool, cmdArr []string, resConf *resource.ResourceConfig) {
	parent, writePipe := container.NewParentProcess(tty)
	if parent == nil {
		return
	}
	if err := parent.Start(); err != nil {
		log.Println(err)
	}

	// create cgroup manager
	// - set resource limit
	// - apply container pid to cgroup
	cgroupManager := cgroups.NewCgroupManager("mydocker-cgroup")
	defer cgroupManager.Destroy()
	_ = cgroupManager.Set(resConf)
	err := cgroupManager.Apply(parent.Process.Pid)
	if err != nil {
		fmt.Println(err)
	}

	sendInitCommand(cmdArr, writePipe)
	_ = parent.Wait()
	container.DeleteWorkSpace("/root/")
	os.Exit(-1)
}

// send init command to child process through pipe
func sendInitCommand(cmdArr []string, writePipe *os.File) {
	commands := strings.Join(cmdArr, " ")
	log.Printf("command all is %s\n", commands)

	writePipe.WriteString(commands)
	writePipe.Close()
}
