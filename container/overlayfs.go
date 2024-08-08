package container

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

// NewWorkSpace creates an OverlayFS as container work space
// NOTE: NewWorkSpace is called in the parent process before the container process(child process) is created
// For now, the container process has not been created yet, so the "container path" is "container path in host"
func NewWorkSpace(rootPath string, volume string) {
	createLower(rootPath)
	createDirs(rootPath)
	mountOverlayFS(rootPath)

	if volume != "" {
		mntPath := path.Join(rootPath, "merged")
		hostPath, containerPath, err := volumeExtract(volume)
		if err != nil {
			fmt.Println(err)
			return
		}

		mountVolume(mntPath, hostPath, containerPath)
	}
}

// createLower uses busybox as the lower of overlayfs
func createLower(rootPath string) {
	busyBoxURL := rootPath + "/busybox/"
	busyBoxTarURL := rootPath + "/busybox.tar"

	// check if busybox dir already exists or not
	fmt.Println(busyBoxURL)
	exist, err := PathExists(busyBoxURL)
	if err != nil {
		fmt.Printf("error while check if busybox exists or not: %v\n", err)
		return
	}

	if !exist {
		if err := os.Mkdir(busyBoxURL, 0777); err != nil {
			fmt.Printf("error while create busybox dir: %v\n", err)
			return
		}
		if _, err := exec.Command("tar", "-xvf", busyBoxTarURL, "-C", busyBoxURL).CombinedOutput(); err != nil {
			fmt.Printf("error while untar busybox.tar: %v\n", err)
			return
		}
	}
}

// createDirs creates upper dir and work dir of overlayfs
func createDirs(rootPath string) {
	upperURL := rootPath + "/upper/"
	workURL := rootPath + "/work/"
	mergedURL := rootPath + "/merged/"

	if err := os.Mkdir(upperURL, 0777); err != nil {
		fmt.Printf("error while create upper dir: %v\n", err)
	}
	if err := os.Mkdir(workURL, 0777); err != nil {
		fmt.Printf("error while create work dir: %v\n", err)
	}
	if err := os.Mkdir(mergedURL, 0777); err != nil {
		fmt.Printf("error while create merged dir: %v\n", err)
	}
}

// mountOverlayFS mounts the overlayfs
func mountOverlayFS(rootPath string) {
	// cmd args: lowerdir=/root/busybox,upperdir=/root/upper,workdir=/root/work
	mountArgs := fmt.Sprintf(
		"lowerdir=%s,upperdir=%s,workdir=%s",
		path.Join(rootPath, "busybox"),
		path.Join(rootPath, "upper"),
		path.Join(rootPath, "work"),
	)

	// mount command:
	// mount -t overlay overlay -o lowerdir=/root/busybox,upperdir=/root/upper,workdir=/root/work /root/merged
	cmd := exec.Command("mount", "-t", "overlay", "overlay", "-o", mountArgs, path.Join(rootPath, "merged"))
	fmt.Printf("mount overlayfs: %v\n", cmd.String())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("error while mount overlayfs: %v", err)
	}
}

// DeleteWorkSpace deletes the overlay filesystem when container exits
func DeleteWorkSpace(rootPath string, volume string) {
	mntPath := path.Join(rootPath, "merged")

	// 如果指定了volume则需要umount volume
	// NOTE: 一定要要先 umount volume ，然后再删除目录，否则由于 bind mount 存在，删除临时目录会导致 volume 目录中的数据丢失。
	if volume != "" {
		_, containerPath, err := volumeExtract(volume)
		if err != nil {
			fmt.Printf("error while extract volume in DeleteWorkSpace(): %v", err)
			return
		}

		umountVolume(mntPath, containerPath)
	}

	umountOverlayFS(path.Join(rootPath, "merged"))
	deleteDirs(rootPath)
}

func umountOverlayFS(mntPath string) {
	cmd := exec.Command("umount", mntPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("error while umount overlayfs: %v\n", err)
	}
}

func deleteDirs(rootPath string) {
	dirs := []string{
		path.Join(rootPath, "upper"),
		path.Join(rootPath, "work"),
		path.Join(rootPath, "merged"),
	}

	for _, dir := range dirs {
		if err := os.RemoveAll(dir); err != nil {
			fmt.Printf("error while remove dir %s: %v", dir, err)
		}
	}
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}