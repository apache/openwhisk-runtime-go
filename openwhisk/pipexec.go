/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package openwhisk

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"
)

// PipeExec is a container of a child process
// connected by input and output.
// The execution model is feeding input lines and and expecting outputs from it line by line.
type PipeExec struct {
	cmd      *exec.Cmd
	scannner *bufio.Scanner
	printer  *bufio.Writer
	stdout   io.Reader
	stderr   io.Reader
	err      error
}

// NewPipeExec creates a child subprocess using the provided command line.
// You can then start it getting a communication channel
func NewPipeExec(command string, args ...string) (proc *PipeExec) {
	cmd := exec.Command(command, args...)
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	pipeOut, pipeIn, _ := os.Pipe()
	cmd.ExtraFiles = []*os.File{pipeIn}
	scanner := bufio.NewScanner(pipeOut)
	printer := bufio.NewWriter(stdin)
	proc = &PipeExec{cmd, scanner, printer, stdout, stderr, nil}
	proc.err = startAndCheck(proc.cmd)
	return
}

// print into the child process
func (proc *PipeExec) print(input string) {
	proc.printer.WriteString(input + "\n")
	proc.printer.Flush()
}

// startAndCheck
func startAndCheck(cmd *exec.Cmd) error {
	//fmt.Println(cmd.Path)
	err := cmd.Start()
	if err != nil {
		return err
	}
	ch := make(chan error)
	go func() { ch <- cmd.Wait() }()
	select {
	case <-ch:
		return fmt.Errorf("command exited")
	case <-time.After(1 * time.Millisecond):
		return nil
	}
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

func service(proc *PipeExec, ch chan string) {
	log.Println("started service")
	for {
		in, ok := <-ch
		if !ok || in == "" {
			log.Println("read failed")
			proc.cmd.Process.Kill()
			if ok {
				log.Println("terminated upon request")
			} else {
				log.Println("terminated: cannot read channel")
			}
			break
		}
		//log.Printf("recv: %s\n", in)
		proc.print(in)
		out := proc.scan()
		if out == "" {
			break
		}
		//log.Printf("sent: %s\n", out)
		ch <- out
	}
	close(ch)
}

func collect(ch chan string, reader io.Reader) {
	scan := bufio.NewScanner(reader)
	for scan.Scan() {
		ch <- scan.Text()
	}
}

func logger(proc *PipeExec, chl chan bool) {

	// poll stdout and stderr
	chOut := make(chan string)
	go collect(chOut, proc.stdout)
	chErr := make(chan string)
	go collect(chErr, proc.stderr)

	// wait for the signal
	for <-chl {
		// flush stdout
		runtime.Gosched()
		for loop := true; loop; {
			select {
			case buf := <-chOut:
				fmt.Println(buf)
			default:
				loop = false
			}
		}
		fmt.Println("XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX")

		// flush stderr
		runtime.Gosched()
		for loop := true; loop; {
			select {
			case buf := <-chErr:
				fmt.Println(buf)
			default:
				loop = false
			}
		}
		fmt.Println("XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX")
	}
	close(chl)
}

// StartService will start a go routine executing a service
func StartService(command string, args ...string) (chan string, chan bool) {
	pipe := NewPipeExec(command, args...)
	if pipe.err != nil {
		log.Print(pipe.err)
		return nil, nil
	}
	// create channel
	ch := make(chan string)
	chl := make(chan bool)

	// read-write loop
	go service(pipe, ch)
	go logger(pipe, chl)
	return ch, chl
}
