package handler

import (
	"log"

	"github.com/454270186/mydocker/container"
	"github.com/spf13/cobra"
)

func InitCmdHandler(cmd *cobra.Command, args []string) {
	log.Println("init start")
	command := args[0]
	log.Printf("command: %s\n", command)

	_ = container.RunContainerInitProcess(command, nil)
}