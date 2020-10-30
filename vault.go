// Copyright 2020 Fairwinds
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vaultutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"k8s.io/klog"
)

// Token is the response of vault token lookup
type Token struct {
	Data struct {
		TTL int `json:"ttl"`
	} `json:"data"`
}

// CheckToken makes sure we have a valid token
func CheckToken() error {
	data, _, err := execute("vault", "token", "lookup", "-format=json")
	if err != nil {
		return err
	}

	token := &Token{}
	err = json.Unmarshal(data, token)
	if err != nil {
		return fmt.Errorf("error unmarshaling vault token: %s", err.Error())
	}

	if token.Data.TTL < 30 {
		return fmt.Errorf("vault token will expire in less than 30 seconds")
	}
	return nil
}

// Login initiates a vault login with the provided method
func Login(loginMethod string) error {
	_, _, err := executeInteractive(true, true, "vault", "login", "-method", loginMethod)
	if err != nil {
		return err
	}
	return nil
}

// revokeLease revokes a vault lease.
// Utilized by individual credential Revoke() functions
func revokeLease(leaseID string) error {
	_, _, err := execute("vault", "lease", "revoke", leaseID)
	if err != nil {
		return err
	}
	return nil
}

// exec returns the output and error of a command run using inventory environment variables.
func execute(name string, arg ...string) ([]byte, string, error) {
	cmd := exec.Command(name, arg...)
	data, err := cmd.CombinedOutput()
	output := strings.TrimSpace(string(data))
	if err != nil {
		return nil, "", fmt.Errorf("exit code %d running command %s: %s", cmd.ProcessState.ExitCode(), cmd.String(), output)
	}
	klog.V(5).Infof("command %s output: %s", cmd.String(), output)
	return data, output, nil
}

// executeInteractive works like exec, but allows interactivity with the command
// https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html
func executeInteractive(showStdErr bool, showStdOut bool, name string, arg ...string) ([]byte, string, error) {

	cmd := exec.Command(name, arg...)

	var stdoutBuf, stderrBuf bytes.Buffer

	if showStdOut {
		cmd.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	}

	if showStdErr {
		cmd.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)
	}

	err := cmd.Run()
	if err != nil {
		return nil, "", fmt.Errorf("cmd.Run() failed with %s", err)
	}
	outStr, errStr := stdoutBuf.String(), stderrBuf.String()
	klog.V(5).Infof("\nstdout:\n%s\nstderr:\n%s\n", outStr, errStr)
	return []byte(outStr), outStr, nil
}
