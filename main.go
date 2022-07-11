package main

import (
	"github.com/gabe565/ruckus-mock-ssh/cli"
	"github.com/gliderlabs/ssh"
	flag "github.com/spf13/pflag"
	"log"
)

func main() {
	var address string
	flag.StringVar(&address, "address", "127.0.0.1:2222", "SSH server listening address")
	var username string
	flag.StringVar(&username, "username", "user", "SSH username")
	var password string
	flag.StringVar(&password, "password", "pass", "SSH password")
	flag.Parse()

	handler := cli.New(username, password)

	log.Println("ssh server running on " + address)
	log.Fatal(ssh.ListenAndServe(address, handler.NewSession))
}
