package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
	"time"
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

func copyAndCapture(w io.Writer, r io.Reader) ([]byte, error) {
	var out []byte
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			out = append(out, d...)
			_, err := w.Write(d)
			if err != nil {
				return out, err
			}
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			return out, err
		}
	}
}

func (vp *VirtualProcess) Read(buf []byte) (readStatus int, err error) {
	readStatus = 0
	err = nil
	ch := make(chan bool)
	go func() {
		readStatus, err = vp.stdout.Read(buf)
		ch <- true
	}()
	select {
	case <-ch:
		return
	case <-time.After(time.Second):
		return 0, errors.New("Timeout")
	}
}

func (vp *VirtualProcess) Wait() {
	vp.proc.Wait()
}
