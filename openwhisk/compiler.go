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
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
)

// check if the file exists and it is already compiled
func isCompiled(file string) bool {
	Debug("IsCompiled? %s", file)
	_, err := os.Stat(file)
	if err != nil {
		return false
	}

	buf, err := ioutil.ReadFile(file)
	if err != nil {
		Debug(err.Error())
		return false
	}

	// check if it is an executable
	return IsExecutable(buf, runtime.GOOS)
}

// CompileAction will compile an anction in source format invoking a compiler
func (ap *ActionProxy) CompileAction(main string, srcDir string, binDir string) error {
	if ap.compiler == "" {
		return fmt.Errorf("No compiler defined")
	}

	Debug("compiling: %s %s %s %s", ap.compiler, main, srcDir, binDir)

	var cmd *exec.Cmd
	cmd = exec.Command(ap.compiler, main, srcDir, binDir)
	cmd.Env = []string{"PATH=" + os.Getenv("PATH")}
	for k, v := range ap.env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}

	// gather stdout and stderr
	out, err := cmd.CombinedOutput()
	Debug("compiler out: %s, %v", out, err)
	if len(out) > 0 {
		return fmt.Errorf("%s", out)
	}
	if err != nil {
		return err
	}
	return nil
}
