package cli

import (
	"github.com/gabe565/ruckus-mock-ssh/cli/cursor"
	"github.com/gliderlabs/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

type Session struct {
	cli                Cli
	session            ssh.Session
	terminal           *terminal.Terminal
	showingCompletions bool
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
		log.Println("auth failure")
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

	term.AutoCompleteCallback = s.AutoComplete

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

func (s *Session) AutoComplete(line string, pos int, key rune) (newLine string, newPos int, ok bool) {
	if s.showingCompletions {
		if _, err := s.terminal.Write([]byte(cursor.ClearScreenBelow())); err != nil {
			log.Println(err)
			return line, pos, false
		}
	}
	s.showingCompletions = false

	switch key {
	case 0x09: // Ctrl-I (Tab)
		result, ok := s.FindCommand(line)
		if ok {
			s.showingCompletions = true
		}
		return result, len(result), ok
	}

	return line, pos, false
}

func (s Session) FindCommand(line string) (completion string, ok bool) {
	var matches []string
	err := fs.WalkDir(responses, "responses", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.Type() == fs.ModeDir {
			return nil
		}

		base := filepath.Base(path)
		base = strings.ReplaceAll(base, "_", " ")
		base = strings.TrimSuffix(base, ".txt")
		if strings.HasPrefix(base, line) {
			matches = append(matches, base)
		}
		return nil
	})
	if err != nil {
		return "", false
	}

	if len(matches) == 1 {
		return matches[0], true
	}

	if len(matches) > 0 {
		matchStr := "\n" + strings.Join(matches, "\n")
		// Move cursor back up to prompt
		matchStr += cursor.MoveUpBeginning(len(matches))
		if _, err := s.terminal.Write([]byte(matchStr)); err != nil {
			log.Println(err)
			return line, false
		}
		return longestCommonPrefix(matches), true
	}

	return line, false
}
