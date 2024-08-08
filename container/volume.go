package container

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/454270186/mydocker/constant"
)

// volumeExtract extracts hostPath and containerPath from given volume by ':'. e.g. /path/in/host:/path/in/container
func volumeExtract(volume string) (pathInHost, pathInContainer string, err error) {
	parts := strings.Split(volume, ":")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid volume [%s], must split by ':'", volume)
	}

	pathInHost, pathInContainer = parts[0], parts[1]
	if pathInHost == "" || pathInContainer == "" {
		return pathInHost, pathInContainer, fmt.Errorf("invalid volume [%s], host path or container path cannot be empty", volume)
	}

	return pathInHost, pathInContainer, nil
}

func mountVolume(mntPath, hostPath, containerPath string) {
	// create dir in host
	if err := os.Mkdir(hostPath, constant.Perm0777); err != nil {
		fmt.Printf("error while create host dir for bind mount: %v\n", err)
	}

	// join the real path of container path in host
	containerPathInHost := path.Join(mntPath, containerPath)
	if err := os.Mkdir(containerPathInHost, constant.Perm0777); err != nil {
		fmt.Printf("error while create container dir for bind mount: %v\n", err)
	}

	// mount -o bind /hostPath /containerPath
	fmt.Println(hostPath, containerPathInHost)
	cmd := exec.Command("mount", "-o", "bind", hostPath, containerPathInHost)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("error while bind mount volume: %v\n", err)
	}
}

func umountVolume(mntPath, containerPath string) {
	// mntPath 为容器在宿主机上的挂载点，例如 /root/merged
	// containerPath 为 volume 在容器中对应的目录，例如 /root/tmp
	// containerPathInHost 则是容器中目录在宿主机上的具体位置，例如 /root/merged/root/tmp
	containerPathInHost := path.Join(mntPath, containerPath)
	cmd := exec.Command("umount", containerPathInHost)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("error while umount volume: %v\n", err)
	}
}