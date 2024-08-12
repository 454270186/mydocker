package handler

import (
	"github.com/454270186/mydocker/container"
	"github.com/spf13/cobra"
)

func PsCmdHandler(cmd *cobra.Command, args []string) {
	container.ListContainers()
}