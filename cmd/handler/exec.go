package handler

import (
	"log"
	"os"

	"github.com/454270186/mydocker/container"
	"github.com/spf13/cobra"
)

func ExecCmdHandler(cmd *cobra.Command, args []string) {
	if os.Getenv(container.EnvExecPid) != "" {
		// if env varible exist, it means that this is the second time calling ExecCmdHandler(), 
		// also means the setns() in Cgo has been executed.
		log.Printf("pid callback pid %v\n", os.Getgid())
		return
	}

	if len(args) < 2 {
		log.Printf("missing container id or command")
		return
	}

	containerId := args[0]
	commandArr := args[1:]

	container.ExecContainer(containerId, commandArr)
}