package cli

import (
	"github.com/gabe565/ruckus-mock-ssh/cli/cursor"
	"github.com/gliderlabs/ssh"
	"golang.org/x/term"
	"io"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

type Session struct {
	cli                Cli
	session            ssh.Session
	terminal           *term.Terminal
	showingCompletions bool
}

func (s Session) Handle() {
	log.Println("begin session")

	t := s.terminal

	defer func() {
		t.SetPrompt("")
		_, _ = t.Write([]byte("Exit ruckus CLI.\n"))
		log.Println("end session")
	}()

	// Login prompt
	t.SetPrompt("\nPlease login: ")
	username, err := t.ReadLine()
	if err != nil {
		return
	}
	log.Println("username: " + username)

	// Password prompt
	password, err := t.ReadPassword("Password: ")
	if err != nil {
		return
	}
	log.Println("password: " + password)

	if !s.cli.AuthPair.Validate(username, password) {
		log.Println("auth failure")
		if _, err := t.Write([]byte("Login incorrect\n")); err != nil {
			log.Println(err)
			return
		}
		return
	}

	log.Println("auth success")

	// Banner
	_, err = t.Write([]byte("Welcome to Ruckus Unleashed Command Line Interface\n"))
	if err != nil {
		log.Println(err)
		return
	}

	t.AutoCompleteCallback = s.AutoComplete

	var superuser bool
	for {
		// Change prompt when superuser
		if superuser {
			t.SetPrompt("ruckus# ")
		} else {
			t.SetPrompt("ruckus> ")
		}

		line, err := t.ReadLine()
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
			if err := s.SendResponse(responseFile); err != nil {
				log.Println(err)
				continue
			}
		}
	}
}

func (s *Session) SendResponse(path string) error {
	path = filepath.Join(
		"responses",
		filepath.Join("/", path+".txt"),
	)

	f, err := responses.Open(path)
	if err != nil {
		return err
	}
	defer func(f fs.File) {
		_ = f.Close()
	}(f)

	if _, err := io.Copy(s.terminal, f); err != nil {
		return err
	}

	return nil
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
		result, ok := s.CompleteCommand(line)
		if ok {
			s.showingCompletions = true
		}
		return result, len(result), ok
	}

	return line, pos, false
}

func (s Session) CompleteCommand(line string) (completion string, ok bool) {
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
