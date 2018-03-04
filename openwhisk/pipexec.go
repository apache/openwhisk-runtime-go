package openwhisk

import (
	"bufio"
	"encoding/json"
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
	err      error
}

// NewPipeExec creates a child subprocess using the provided command line.
// You can then start it getting a communcation channel
func NewPipeExec(command string, args ...string) (proc *PipeExec) {
	cmd := exec.Command(command, args...)
	in, _ := cmd.StdinPipe()
	out, _ := cmd.StdoutPipe()
	scanner := bufio.NewScanner(out)
	printer := bufio.NewWriter(in)
	proc = &PipeExec{cmd, scanner, printer, nil}
	proc.err = proc.cmd.Start()
	proc.handshake()
	log.Println(proc.err)
	return
}

// print into the child process
func (proc *PipeExec) print(input string) {
	proc.printer.WriteString(input + "\n")
	proc.printer.Flush()
}

// scan from the child process
func (proc *PipeExec) scan() string {
	if proc.err == nil {
		if proc.scannner.Scan() {
			return proc.scannner.Text()
		}
		proc.err = proc.scannner.Err()
	}
	return ""
}

// handshake
func (proc *PipeExec) handshake() {
	var welcome struct {
		OpenWhisk int
	}
	if proc.scannner.Scan() {
		buf := proc.scannner.Bytes()
		proc.err = json.Unmarshal(buf, &welcome)
		if proc.err != nil {
			return
		}
		if welcome.OpenWhisk < 1 {
			proc.err = fmt.Errorf("failed handshake: %s", string(buf))
		}
	} else {
		proc.err = fmt.Errorf("no handshake")
	}
}

func service(proc *PipeExec, ch chan string) {
	log.Println("started service")
	for {
		in, ok := <-ch
		if !ok || in == "" {
			proc.cmd.Process.Kill()
			log.Println("terminated")
			break
		}
		log.Printf("recv: %s\n", in)
		proc.print(in)
		out := proc.scan()
		if out == "" {
			break
		}
		log.Printf("sent: %s\n", out)
		ch <- out
	}
	close(ch)
}

// StartService will start a go routine executing a service
func StartService(command string, args ...string) chan string {
	pipe := NewPipeExec(command, args...)
	if pipe.err != nil {
		log.Print(pipe.err)
		return nil
	}
	ch := make(chan string)
	go service(pipe, ch)
	return ch
}
