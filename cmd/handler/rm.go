package handler

import (
	"log"

	"github.com/454270186/mydocker/container"
	"github.com/spf13/cobra"
)

var (
	IsForce bool
)

func RmCmdHandler(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		log.Println("missing container id")
		return
	}

	containerId := args[0]
	err := container.RemoveContainer(containerId, IsForce)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Successfully remove container %s\n", containerId)
}