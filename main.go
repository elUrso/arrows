package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/pkg/term"
)

func main() {
	// get tty file path using the tty command

	cmd := exec.Command("tty")

	// input must be from current tty, otherwise it fails
	cmd.Stdin = os.Stdin

	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}

	// trim the new line
	termpath := strings.TrimSpace(string(out))

	// open the terminal control file
	t, err := term.Open(termpath)
	if err != nil {
		log.Fatal(err)
	}

	// defer the restore to return it to the normal state
	defer t.Restore()

	// restore also in case of a control-c
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		t.Restore()
		os.Exit(0)
	}()

	// set to cbreak, uncooked mode
	t.SetCbreak()

	// buffer to read the inputs
	b := make([]byte, 3)

	// q for quit
	for b[0] != 'q' {
		os.Stdin.Read(b)

		// arrows are marked by the sequence:
		// esc (27) [ (91) and a letter
		if b[0] == 27 && b[1] == 91 {
			switch b[2] {
			case 65: // A
				fmt.Print("⬆️ ")
			case 67: // C
				fmt.Print("➡️ ")
			case 66: // B
				fmt.Print("⬇️ ")
			case 68: // D
				fmt.Print("⬅️ ")
			}
		}
	}
}
