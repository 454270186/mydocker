package handler

import (
	"log"
	"os"
	"strings"

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

	Run(IsTTY, args)
}

func Run(tty bool, cmdArr []string) {
	parent, writePipe := container.NewParentProcess(tty)
	if parent == nil {
		return
	}
	if err := parent.Start(); err != nil {
		log.Println(err)
	}
	
	sendInitCommand(cmdArr, writePipe)
	_ = parent.Wait()
	os.Exit(-1)
}

// send init command to child process through pipe
func sendInitCommand(cmdArr []string, writePipe *os.File) {
	commands := strings.Join(cmdArr, " ")
	log.Printf("command all is %s\n", commands)
	
	writePipe.WriteString(commands)
	writePipe.Close()
}