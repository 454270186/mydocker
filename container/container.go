package container

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

func NewParentProcess(tty bool) (*exec.Cmd, *os.File) {
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
	return cmd, writePipe
}

func RunContainerInitProcess(commandArr []string) error {
	mountProc()

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

func mountProc() {
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	_ = syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
}