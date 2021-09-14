/*
Copyright 2015 The Kubernetes Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

*/

package util

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"k8s.io/klog/v2"
)

const (
	ShellToUse            = "bash"
	CommandTimeoutSeconds = 10
)

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func ExecCommandLine(commandline string) (string, error) {
	return ExecCommandLineWithTimeOut(commandline, CommandTimeoutSeconds)
}

func ExecCommandLineWithTimeOut(commandline string, timeout int) (string, error) {
	var cmd *exec.Cmd

	/*if !strings.Contains(commandline, " mission") && !strings.Contains(commandline, "cluster-info") && !strings.Contains(commandline, "get nodes") &&  !strings.Contains(commandline, "get edgeclusters") {
		klog.Infof("~~~~~~~~ 1 running (%v) (%v)", commandline, timeout)
	}*/

	background := false
	if strings.HasSuffix(commandline, " &") {
		background = true
		//commandline = commandline[0: len(commandline)-2]
		timeout = 0
	}

	if timeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
		defer cancel()

		cmd = exec.CommandContext(ctx, ShellToUse, "-c", commandline)
	} else {
		cmd = exec.Command(ShellToUse, "-c", commandline)
	}

	/*if !strings.Contains(commandline, " mission") && !strings.Contains(commandline, "cluster-info") && !strings.Contains(commandline, "get nodes") &&  !strings.Contains(commandline, "get edgeclusters") {
		klog.Infof("~~~~~~~~ 2 running (%v) (%v)", commandline, timeout)
	}*/
	exitCode := 0
	var output []byte
	var err error

	if background {
		fmt.Printf("\n Qian 1 --------- %v", commandline)
		if err := cmd.Start(); err != nil {
			fmt.Printf("\n Qian 2 --------- %v", commandline)
			output = []byte("Failed to start the command")
			exitCode = 1
		}
		fmt.Printf("\n Qian 3 --------- %v", commandline)
		//go cmd.CombinedOutput()	
		//go cmd.CombinedOutput()
	} else {
		if output, err = cmd.CombinedOutput(); err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				exitCode = exitError.ExitCode()
			}
		}
	}

	/*if !strings.Contains(commandline, " mission") && !strings.Contains(commandline, "cluster-info") && !strings.Contains(commandline, "get nodes") &&  !strings.Contains(commandline, "get edgeclusters") {
		klog.Infof("~~~~~~~~ 3 running (%v) (%v)", commandline, timeout)
	}*/
	var finalErr error
	if exitCode != 0 || err != nil {
		finalErr = fmt.Errorf("Command (%v) failed: exitcode: %v, output (%v), error: %v", commandline, exitCode, string(output), err)
	} else {
		klog.V(3).Infof("Running Command (%v) succeeded", commandline)
	}

	/*if !strings.Contains(commandline, " mission") && !strings.Contains(commandline, "cluster-info") && !strings.Contains(commandline, "get nodes") &&  !strings.Contains(commandline, "get edgeclusters") {
		klog.Infof("~~~~~~~~ 4 running (%v) (%v)", commandline, timeout)
	}*/
	return string(output), finalErr
}
