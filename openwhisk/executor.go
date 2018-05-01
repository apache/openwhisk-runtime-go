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
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"
)

// TIMEOUT to wait for process to start
// and log to be produced
const TIMEOUT = 5 * time.Millisecond

// Executor is the container and the guardian  of a child process
// It starts a command, feeds input and output, read logs and control its termination
type Executor struct {
	io      chan string
	log     chan bool
	exit    chan error
	_cmd    *exec.Cmd
	_input  *bufio.Writer
	_output *bufio.Scanner
	_logout *bufio.Scanner
	_logerr *bufio.Scanner
	_logbuf *os.File
}

// NewExecutor creates a child subprocess using the provided command line,
// writing the logs in the given file.
// You can then start it getting a communication channel
func NewExecutor(logbuf *os.File, command string, args ...string) (proc *Executor) {
	cmd := exec.Command(command, args...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil
	}

	pipeOut, pipeIn, err := os.Pipe()
	if err != nil {
		return nil
	}
	cmd.ExtraFiles = []*os.File{pipeIn}

	return &Executor{
		make(chan string),
		make(chan bool),
		make(chan error),
		cmd,
		bufio.NewWriter(stdin),
		bufio.NewScanner(pipeOut),
		bufio.NewScanner(stdout),
		bufio.NewScanner(stderr),
		logbuf,
	}
}

// collect log from a stream
func _collect(ch chan string, scan *bufio.Scanner) {
	for scan.Scan() {
		ch <- scan.Text()
	}
}

// loop over the command executing
// returning when the command exits
func (proc *Executor) run() {
	log.Println("run: start")
	err := proc._cmd.Start()
	if err != nil {
		proc.exit <- err
		log.Println("run: early exit")
		proc._cmd = nil // do not kill
		return
	}
	// wait for the exit
	proc.exit <- proc._cmd.Wait()
	proc._cmd = nil // do not kill
	log.Println("run: end")
}

func (proc *Executor) drain(ch chan string) {
	runtime.Gosched()
	for loop := true; loop; {
		select {
		case buf := <-ch:
			fmt.Fprintln(proc._logbuf, buf)
		case <-time.After(TIMEOUT):
			loop = false
		}
	}
	fmt.Fprintln(proc._logbuf, "XXX_THE_END_OF_A_WHISK_ACTIVATION_XXX")
}

// manage copying stdout and stder in output
// with log guards
func (proc *Executor) logger() {
	log.Println("logger: start")
	// poll stdout and stderr
	chOut := make(chan string)
	go _collect(chOut, proc._logout)
	chErr := make(chan string)
	go _collect(chErr, proc._logerr)

	// wait for the signal
	for <-proc.log {
		// flush stdout
		proc.drain(chOut)
		// flush stderr
		proc.drain(chErr)
	}
	log.Printf("logger: end")
}

// main service function
// writing in input
// and reading in output
// using the provide channels
func (proc *Executor) service() {
	log.Println("service: start")
	for {
		in := <-proc.io
		if in == "" {
			log.Println("terminated upon request")
			break
		}
		// input/output with the process
		log.Printf(">>>%s\n", in)
		proc._input.WriteString(in + "\n")
		proc._input.Flush()
		if proc._output.Scan() {
			out := proc._output.Text()
			log.Printf("<<<%s\n", out)
			proc.io <- out
			if out == "" {
				break
			}
		}
	}
	log.Printf("service: end")
}

// Start execution of the command
// returns an error if the program fails
func (proc *Executor) Start() error {
	// start the underlying executable
	// check if died
	go proc.run()
	select {
	case <-proc.exit:
		// oops, it died
		return fmt.Errorf("command exited")
	case <-time.After(TIMEOUT):
		// ok let's process it
		go proc.service()
		go proc.logger()
	}
	return nil
}

// Stop will kill the process
// and close the channels
func (proc *Executor) Stop() {
	log.Println("stopping")
	if proc._cmd != nil {
		proc.log <- false
		proc.io <- ""
		proc._cmd.Process.Kill()
		<-proc.exit
		proc._cmd = nil
	}
	close(proc.io)
	close(proc.exit)
	close(proc.log)
}
