package cli

import (
	"io"
)

// LineToCarriage Converts input `\n` to `\r` to fix issues with
// the Python pexpect library and ssh.Terminal
type LineToCarriage struct {
	rw io.ReadWriter
}

func (r LineToCarriage) Read(p []byte) (n int, err error) {
	if n, err = r.rw.Read(p); err != nil {
		return n, err
	}
	for k, v := range p {
		if v == '\n' {
			p[k] = '\r'
		}
	}
	return n, err
}

func (r LineToCarriage) Write(p []byte) (n int, err error) {
	return r.rw.Write(p)
}
