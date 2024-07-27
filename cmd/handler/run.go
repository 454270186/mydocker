package handler

import (
	"log"
	"os"

	"github.com/454270186/mydocker/container"
	"github.com/spf13/cobra"
)

var (
	IsTTY bool
)

func RunCmdHandler(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		log.Println("missing container command")
		return
	}

	command := args[0]
	Run(IsTTY, command)
}

func Run(tty bool, cmd string) {
	parent := container.NewParentProcess(tty, cmd)
	if err := parent.Start(); err != nil {
		log.Println(err)
	}
	
	_ = parent.Wait()
	os.Exit(-1)
}