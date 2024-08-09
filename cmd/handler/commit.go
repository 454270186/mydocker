package handler

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/spf13/cobra"
)

func CommitCmdHandler(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		log.Println("missing image name")
		return
	}

	imageName := args[0]
	commitContainer(imageName)
}

// commitContainer tar container's rootfs into image.tar
func commitContainer(imageName string) {
	mntPath := "/root/merged"
	imageTar := "/root/" + imageName + ".tar"
	fmt.Println("commitContainer imageTar:", imageTar)

	if _, err := exec.Command("tar", "-czf", imageTar, "-C", mntPath, ".").CombinedOutput(); err != nil {
		fmt.Printf("error while tar dir %s: %v\n", imageTar, err)
	}
}
