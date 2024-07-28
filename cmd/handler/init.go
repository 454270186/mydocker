package handler

import (
	"io"
	"log"
	"os"
	"strings"

	"github.com/454270186/mydocker/container"
	"github.com/spf13/cobra"
)

const (
	readPipeFdIndex = 3
)

func InitCmdHandler(cmd *cobra.Command, args []string) {
	log.Println("init start")

	// read commands from pipe
	commandArr := readUserCommand()
	if len(commandArr) == 0 {
		log.Println("user command is empty")
		return
	}

	_ = container.RunContainerInitProcess(commandArr)
}

func readUserCommand() []string {
	pipe := os.NewFile(uintptr(readPipeFdIndex), "pipe")
	msg, err := io.ReadAll(pipe)
	if err != nil {
		log.Println(err)
		return nil
	}

	return strings.Split(string(msg), " ")
}