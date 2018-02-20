package openwhisk

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
)

// PipeExec is a container of a child process
// connected by input and output.
// The execution model is feeding input lines and and expecting outputs from it line by line.
type PipeExec struct {
	cmd      *exec.Cmd
	scannner *bufio.Scanner
	printer  *bufio.Writer
}

// NewPipeExec creates a child subprocess using the provided command line.
// You can then start it getting a communcation channel
func NewPipeExec(command string, args ...string) (proc *PipeExec) {
	cmd := exec.Command(command, args...)
	in, _ := cmd.StdinPipe()
	out, _ := cmd.StdoutPipe()
	scanner := bufio.NewScanner(out)
	printer := bufio.NewWriter(in)
	proc = &PipeExec{cmd, scanner, printer}
	proc.cmd.Start()
	return
}

// print into the child process
func (proc *PipeExec) print(input string) {
	proc.printer.WriteString(input + "\n")
	proc.printer.Flush()
}

// scan from the child process
func (proc *PipeExec) scan() string {
	if proc.scannner.Scan() {
		return proc.scannner.Text()
	}
	return ""
}

func service(proc *PipeExec, ch chan string) {
	log.Println("started service")
	for {
		in := <-ch
		log.Printf("recv: %s\n", in)
		if in == "" {
			// TODO: test this
			proc.cmd.Process.Kill()
			fmt.Println("terminated")
			break
		}
		proc.print(in)
		out := proc.scan()
		log.Printf("sent: %s\n", out)
		ch <- out
	}
}

// StartService will start a go routine executing a service
func StartService(command string, args ...string) chan string {
	pipe := NewPipeExec(command, args...)
	ch := make(chan string)
	go service(pipe, ch)
	return ch
}
