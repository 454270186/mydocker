package handler

import (
	"log"

	"github.com/454270186/mydocker/container"
	"github.com/spf13/cobra"
)

func StopCmdHander(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		log.Println("missing container id")
		return
	}

	containerId := args[0]
	err := container.StopContainer(containerId)
	if err != nil {
		log.Println(err)
	}
}