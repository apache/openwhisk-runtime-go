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

// IsElf checks for a Linux executable
func IsElf(buf []byte) bool {
	return len(buf) > 52 &&
		buf[0] == 0x7F && buf[1] == 0x45 &&
		buf[2] == 0x4C && buf[3] == 0x46
}

// IsExe checks for a Windows executable
func IsExe(buf []byte) bool {
	return len(buf) > 1 &&
		buf[0] == 0x4D && buf[1] == 0x5A
}

// IsMach64 checks for OSX executable
func IsMach64(buf []byte) bool {
	return len(buf) > 4 &&
		buf[0] == 0xcf && buf[1] == 0xfa &&
		buf[2] == 0xed && buf[3] == 0xfe
}

// IsBangPath checks for a shell executable
func IsBangPath(buf []byte) bool {
	return len(buf) > 2 &&
		buf[0] == '#' && buf[1] == '!'
}

// IsExecutable check if it is an executable, according the current runtime
func IsExecutable(buf []byte, runtime string) bool {
	Debug("checking executable for %s", runtime)
	switch runtime {
	case "darwin":
		return IsMach64(buf) || IsBangPath(buf)
	case "linux":
		return IsElf(buf) || IsBangPath(buf)
	case "windows":
		return IsExe(buf)
	default:
		return false
	}
}

// IsZip checks if it is a zip file
func IsZip(buf []byte) bool {
	return len(buf) > 3 &&
		buf[0] == 0x50 && buf[1] == 0x4B &&
		(buf[2] == 0x3 || buf[2] == 0x5 || buf[2] == 0x7) &&
		(buf[3] == 0x4 || buf[3] == 0x6 || buf[3] == 0x8)
}

// IsGz checks if the given file is a valid tar.gz file
func IsGz(buf []byte) bool {
	// Magic number: The first two bytes are fixed (0x1f and 0x8b), which represent the magic number of a gzip file.
	// Compression method: The third byte indicates the compression method used. For gzip files, this is always  (deflate).
	return len(buf) > 3 &&
		buf[0] == 0x1f &&
		buf[1] == 0x8b &&
		buf[2] == 0x08
}
