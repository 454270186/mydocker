package handler

import (
	"log"

	"github.com/454270186/mydocker/container"
	"github.com/spf13/cobra"
)

func LogsCmdHandler(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		log.Println("missing container name")
		return
	}

	containerId := args[0]
	container.LogContainer(containerId)
}