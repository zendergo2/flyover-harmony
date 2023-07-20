package main

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

type VirtualProcess struct {
	proc   *os.Process
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr bytes.Buffer
}

func (vp *VirtualProcess) Start() {
	// convert to StartProcess instead of command
	wd, _ := os.Getwd()
	cmd := exec.Command(wd + "/test.sh")
	vp.stdin, _ = cmd.StdinPipe()
	vp.stdout, _ = cmd.StdoutPipe()
	cmd.Start()
}

func (vp *VirtualProcess) Wait() {
	vp.proc.Wait()
}
