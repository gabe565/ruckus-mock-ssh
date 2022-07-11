package cli

import (
	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

func New(username, password string) Cli {
	return Cli{
		AuthPair: AuthPair{
			Username: username,
			Password: password,
		},
	}
}

type Cli struct {
	AuthPair AuthPair
}

func (cli Cli) NewSession(sshSession ssh.Session) {
	session := Session{
		cli:     cli,
		session: sshSession,
		terminal: terminal.NewTerminal(
			LineToCarriage{sshSession},
			"",
		),
	}
	session.Handle()
}
