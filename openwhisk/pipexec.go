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
	logger   *bufio.Scanner
	err      error
}

// NewPipeExec creates a child subprocess using the provided command line.
// You can then start it getting a communcation channel
func NewPipeExec(command string, args ...string) (proc *PipeExec) {
	cmd := exec.Command(command, args...)
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	scanner := bufio.NewScanner(stdout)
	printer := bufio.NewWriter(stdin)
	logger := bufio.NewScanner(stderr)

	proc = &PipeExec{cmd, scanner, printer, logger, nil}
	proc.err = proc.cmd.Start()
	if proc.err == nil {
		proc.handshake()
	}
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

// handshake - leave the error in proc.err
func (proc *PipeExec) handshake() {
	var welcome struct {
		OpenWhisk int
	}
	if proc.scannner.Scan() {
		buf := proc.scannner.Bytes()
		proc.err = json.Unmarshal(buf, &welcome)
		if proc.err != nil || welcome.OpenWhisk < 1 {
			proc.err = fmt.Errorf("failed handshake: %s", string(buf))
		}
	} else {
		proc.err = fmt.Errorf("no handshake")
	}
}

func logger(proc *PipeExec) {
	log.Println("started logger")
	// scanner read stderr continuosly
	// it will exit when the underlying process terminate
	for proc.logger.Scan() {
		fmt.Println(proc.logger.Text())
	}
	log.Println("exited logger")
}

func service(proc *PipeExec, ch chan string) {
	log.Println("started service")
	for {
		in, ok := <-ch
		if !ok || in == "" {
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

// StartService will start a go routine executing a service
func StartService(command string, args ...string) chan string {
	pipe := NewPipeExec(command, args...)
	if pipe.err != nil {
		log.Print(pipe.err)
		return nil
	}
	// create channer
	ch := make(chan string)
	// read-write loop
	go service(pipe, ch)
	// stderr to stdout
	go logger(pipe)
	return ch
}
