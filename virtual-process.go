package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
)

type VirtualProcess struct {
	proc   *os.Process
	stdin  bytes.Buffer
	stdout bytes.Buffer
	stderr bytes.Buffer
}

func (vp *VirtualProcess) Start() {
	wd, _ := os.Getwd()
	cmd := exec.Command(wd + "/test.sh")
	cmd.Stdin = &vp.stdin
	cmd.Stdout = &vp.stdout
	cmd.Stderr = &vp.stderr
	cmd.Run()
	log.Println("stdout:", vp.stdout.String())
	if vp.stderr.Len() > 0 {
		log.Println("stderr:", vp.stderr.String())
	}
	var out = vp.stdout.Bytes()
	if bytes.HasSuffix(out, []byte("\n")) {
		vp.stdout.Truncate(len(out) - 1)
		vp.stdout.Write([]byte{'\r', '\n'})
	}
}

func (vp *VirtualProcess) Wait() {
	//vp.proc.Wait()
}
