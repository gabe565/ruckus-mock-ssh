package cli

import (
	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

type Session struct {
	cli      Cli
	session  ssh.Session
	terminal *terminal.Terminal
}

func (s Session) Handle() {
	log.Println("begin session")
	defer log.Println("end session")

	term := s.terminal

	// Login prompt
	term.SetPrompt("Please login: ")
	username, err := term.ReadLine()
	if err != nil {
		return
	}
	log.Println("username: " + username)

	// Password prompt
	password, err := term.ReadPassword("Password: ")
	if err != nil {
		return
	}
	log.Println("password: " + password)

	if !s.cli.AuthPair.Validate(username, password) {
		if _, err := term.Write([]byte("Login incorrect\n")); err != nil {
			log.Println(err)
			return
		}
		return
	}

	log.Println("auth success")

	// Banner
	_, err = term.Write([]byte("Welcome to Ruckus Unleashed Command Line Interface\n"))
	if err != nil {
		log.Println(err)
		return
	}

	var superuser bool
	for {
		// Change prompt when superuser
		if superuser {
			term.SetPrompt("ruckus# ")
		} else {
			term.SetPrompt("ruckus> ")
		}

		line, err := term.ReadLine()
		if err != nil {
			log.Println(err)
			return
		}

		log.Println("input: " + line)
		split := strings.Split(line, " ")

		responseFile := strings.ReplaceAll(line, " ", "_")
		switch split[0] {
		case "quit", "exit", "abort", "end", "bye":
			return
		case "?":
			responseFile = "help"
		case "enable":
			superuser = true
			continue
		case "disable":
			superuser = false
			continue
		}

		if responseFile != "" {
			responseFile = filepath.Join(
				"responses",
				filepath.Join("/", responseFile+".txt"),
			)
			response, err := fs.ReadFile(responses, responseFile)
			if err != nil {
				log.Println(err)
				continue
			}

			if _, err := term.Write(response); err != nil {
				log.Println(err)
				return
			}
		}
	}
}
