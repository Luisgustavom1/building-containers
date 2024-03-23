// go:build linux
//go:build linux
// +build linux

package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	switch os.Args[1] {
	case "run":
		parent()
	case "child":
		child()
	default:
		panic("not found command")
	}
}

func parent() {
	fmt.Println("running parent process with pid: ", os.Getpid())
	childArgs := append([]string{"child"}, os.Args[2:]...)
	// /proc/self is a real symbolic link to the /proc/ subdirectory of the process that is making the call.
	// https://elixir.bootlin.com/linux/latest/source/fs/proc/self.c
	cmd := exec.Command("/proc/self/exe", childArgs...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}
}

func child() {
	fmt.Println("running child process with pid: ", os.Getpid())

	syscall.Sethostname([]byte("my-container"))
	syscall.Chroot("/home/Ubuntu/building-containers")
	syscall.Chdir("/")
	syscall.Mount("proc", "proc", "proc", 0, "")

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("error: ", err)
		os.Exit(1)
	}

	syscall.Unmount("/proc", 0)
}