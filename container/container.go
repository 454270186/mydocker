package container

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"syscall"
)

func NewParentProcess(tty bool, volume string) (*exec.Cmd, *os.File) {
	// create pipe
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		log.Println(err)
		return nil, nil
	}

	args := []string{"init"}
	cmd := exec.Command("/proc/self/exe", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}

	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	// pass the read pipe to child process
	cmd.ExtraFiles = []*os.File{readPipe}

	// init overlayfs workspace
	rootPath := "/root"
	NewWorkSpace(rootPath, volume)
	cmd.Dir = path.Join(rootPath, "merged")

	return cmd, writePipe
}

func RunContainerInitProcess(commandArr []string) error {
	setMount()

	path, err := exec.LookPath(commandArr[0])
	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("find path %s\n", path)
	if err := syscall.Exec(path, commandArr[0:], os.Environ()); err != nil {
		log.Println(err)
	}

	return nil
}

// setMount init mount point
func setMount() {
	pwd, err := os.Getwd() // pwd(work dir) is set when this child process was created
	if err != nil {
		fmt.Printf("error while get pwd: %v", err)
		return
	}

	fmt.Printf("Current location is %s\n", pwd)

	// 显示声明当前这个新的mount namespace独立
	err = syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	if err != nil {
		fmt.Printf("error while private mount: %v\n", err)
		return
	}

	err = pivotRoot(pwd)
	if err != nil {
		fmt.Println(err)
		return
	}

	// mount /proc
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	_ = syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	// 由于前面 pivotRoot 切换了 rootfs，因此这里重新 mount 一下 /dev 目录
	// tmpfs 是基于 件系 使用 RAM、swap 分区来存储。
	// 不挂载 /dev，会导致容器内部无法访问和使用许多设备，这可能导致系统无法正常工作
	syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
}

func pivotRoot(root string) error {
	cur, _ := os.Getwd()
	fmt.Println(cur) // output: /root/busybox

	err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, "")
	if err != nil {
		return fmt.Errorf("error while mounting rootfs to itself: %v", err)
	}

	// 创建 rootfs/.pivot_root 目录用于存储 old_root
	pivotDir := filepath.Join(root, ".pivot_root")
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		return fmt.Errorf("error while make old root dir: %v", err)
	}

	// 执行pivot_root调用,将系统rootfs切换到新的rootfs,
	// PivotRoot调用会把 old_root挂载到pivotDir,也就是rootfs/.pivot_root,挂载点现在依然可以在mount命令中看到
	err = syscall.PivotRoot(root, pivotDir)
	if err != nil {
		return fmt.Errorf("error while pivotRoot, new_root %s, put_old %s: %v", root, pivotDir, err)
	}

	cur, _ = os.Getwd()
	fmt.Println(cur) // output: /.pivot_root/root/busybox

	// 修改当前的工作目录到根目录
	err = syscall.Chdir("/")
	if err != nil {
		return fmt.Errorf("error while chdir to /: %v", err)
	}

	cur, _ = os.Getwd()
	fmt.Println(cur) // output: /

	// 最后再把old_root umount了，即 umount rootfs/.pivot_root
	// 由于当前已经是在 rootfs 下了，就不能再用上面的rootfs/.pivot_root这个路径了,现在直接用/.pivot_root这个路径即可
	pivotDir = filepath.Join("/", ".pivot_root")
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("error while unmount pivot_root dir: %v", err)
	}

	if err := os.Remove(pivotDir); err != nil {
		return err
	}

	return nil
}