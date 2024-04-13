package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	if len(os.Args) == 1 || os.Args[1] != "child" {
		runContainer()
	} else {
		initContainer()
	}
}

func runContainer() {
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[1:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	if err := cmd.Run(); err != nil {
		fmt.Println("Error running container:", err)
		return
	}
}

func initContainer() {
	// 设置容器的namespace和cgroup
	must(syscall.Sethostname([]byte("container")))
	must(syscall.Chroot("/"))
	must(syscall.Chdir("/"))
	must(syscall.Mount("proc", "proc", "proc", 0, ""))

	// 在容器内打印一句话
	fmt.Println("Hello from the container!")

	// 在容器内启动一个交互式shell
	cmd := exec.Command("/bin/sh")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	must(cmd.Run())

	must(syscall.Unmount("proc", 0))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
