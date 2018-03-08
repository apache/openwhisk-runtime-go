package openwhisk

import (
	"fmt"
	"os/exec"
	"syscall"
)

func start(command string) {
	fmt.Println("-- start --")
	cmd := exec.Command(command)
	err := cmd.Start()
	if err == nil {
		//fmt.Println(cmd.ProcessState.Exited())
		//_, err := os.FindProcess(cmd.Process.Pid)
		//fmt.Printf("find err=%v\n", err)
		//fmt.Printf("start %s pid %d err=%v\n", command, cmd.Process.Pid, err)
		err := cmd.Process.Signal(syscall.Signal(0))
		if err != nil {
			fmt.Printf("cannot find %s err=%v\n", command, err)
		} else {
			fmt.Printf("all ok %s\n", command)
		}
	} else {
		fmt.Printf("cannot start %s err=%v\n", command, err)
	}
}

func Example_demo() {
	start("donotexist")
	start("/etc/passwd")
	start("true")
	start("bc")
	return
	// Output:
	// test
}
