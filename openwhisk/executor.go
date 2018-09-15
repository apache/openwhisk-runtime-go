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
	"os"
	"os/exec"
	"runtime"
	"time"
)

// DefaultTimeoutInit to wait for a process to start
var DefaultTimeoutInit = 5 * time.Millisecond

// DefaultTimeoutDrain to wait for draining logs
var DefaultTimeoutDrain = 5 * time.Millisecond

// Executor is the container and the guardian  of a child process
// It starts a command, feeds input and output, read logs and control its termination
type Executor struct {
	io      chan []byte
	log     chan bool
	exit    chan error
	_cmd    *exec.Cmd
	_input  io.WriteCloser
	_output *bufio.Reader
	_logout *bufio.Reader
	_logerr *bufio.Reader
	_outbuf *os.File
	_errbuf *os.File
}

// NewExecutor creates a child subprocess using the provided command line,
// writing the logs in the given file.
// You can then start it getting a communication channel
func NewExecutor(outbuf *os.File, errbuf *os.File, command string, args ...string) (proc *Executor) {
	cmd := exec.Command(command, args...)
	cmd.Env = []string{
		"__OW_API_HOST=" + os.Getenv("__OW_API_HOST"),
	}
	if Debugging {
		cmd.Env = append(cmd.Env, "OW_DEBUG=/tmp/action.log")
	}
	Debug("env: %v", cmd.Env)

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

	pout := bufio.NewReader(pipeOut)
	sout := bufio.NewReader(stdout)
	serr := bufio.NewReader(stderr)

	return &Executor{
		make(chan []byte),
		make(chan bool),
		make(chan error),
		cmd,
		stdin,
		pout,
		sout,
		serr,
		outbuf,
		errbuf,
	}
}

// collect log from a stream
func _collect(ch chan string, reader *bufio.Reader) {
	for {
		buf, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}
		ch <- string(buf)
	}
}

// loop over the command executing
// returning when the command exits
func (proc *Executor) run() {
	Debug("run: start")
	err := proc._cmd.Start()
	if err != nil {
		proc.exit <- err
		Debug("run: early exit")
		proc._cmd = nil // do not kill
		return
	}
	Debug("pid: %d", proc._cmd.Process.Pid)
	// wait for the exit
	proc.exit <- proc._cmd.Wait()
	proc._cmd = nil // do not kill
	Debug("run: end")
}

func drain(ch chan string, out *os.File) {
	for loop := true; loop; {
		runtime.Gosched()
		select {
		case buf := <-ch:
			fmt.Fprint(out, buf)
		case <-time.After(DefaultTimeoutDrain):
			loop = false
		}
	}
	fmt.Fprintln(out, OutputGuard)
}

// manage copying stdout and stder in output
// with log guards
func (proc *Executor) logger() {
	Debug("logger: start")
	// poll stdout and stderr
	chOut := make(chan string)
	go _collect(chOut, proc._logout)
	chErr := make(chan string)
	go _collect(chErr, proc._logerr)

	// wait for the signal
	for <-proc.log {
		// flush stdout
		drain(chOut, proc._outbuf)
		// flush stderr
		drain(chErr, proc._errbuf)
	}
	proc._outbuf.Sync()
	proc._errbuf.Sync()
	Debug("logger: end")
}

// main service function
// writing in input
// and reading in output
// using the provide channels
func (proc *Executor) service() {
	Debug("service: start")
	for {
		in := <-proc.io
		if len(in) == 0 {
			Debug("terminated upon request")
			break
		}
		// input to the subprocess
		DebugLimit(">>>", in, 120)
		proc._input.Write(in)
		proc._input.Write([]byte("\n"))
		Debug("done")

		// ok now give a chance to run to goroutines
		runtime.Gosched()

		// input to the subprocess
		out, err := proc._output.ReadBytes('\n')
		if err != nil {
			break
		}
		DebugLimit("<<<", out, 120)
		proc.io <- out
		if len(out) == 0 {
			Debug("empty input - exiting")
			break
		}
	}
	Debug("service: end")
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
	case <-time.After(DefaultTimeoutInit):
		// ok let's process it
		go proc.service()
		go proc.logger()
	}
	return nil
}

// Stop will kill the process
// and close the channels
func (proc *Executor) Stop() {
	Debug("stopping")
	if proc._cmd != nil {
		proc.log <- false
		proc.io <- []byte("")
		proc._cmd.Process.Kill()
		<-proc.exit
		proc._cmd = nil
	}
	close(proc.io)
	close(proc.exit)
	close(proc.log)
}
