package main

import (
	"fmt"
	"github.com/inancgumus/screen"
	"os"
	"os/exec"
	"strings"
)

func parse(line []string) int {

	for i, s := range line {
		line[i] = strings.TrimSpace(s)
	}

	switch line[0] {
	case "cd":
		homedir, err := os.UserHomeDir()

		if len(line) == 1 {
			err = os.Chdir(homedir)
			if err != nil {
				return 0
			}
		} else {
			err = os.Chdir(line[1])
			if err != nil {
				return 0
			}
		}
	case "pwd":
		getwd, err := os.Getwd()
		if err != nil {
			return 0
		}
		println(getwd)
	case "exit":
		err := os.Remove(socketPath)
		if err != nil {
			os.Exit(55)
		}
		os.Exit(0)
	case "ls":
		if len(line) == 1 {
			getwd, err := os.Getwd()
			if err != nil {
				return 0
			}
			line = append(line, getwd)
		}
		dir, err := os.ReadDir(line[1])

		if err != nil {
			return 0
		}

		for _, entry := range dir {
			println(entry.Name())
		}
	case "clear":
		screen.Clear()
		screen.MoveTopLeft()
	case "send_sock":
		var sendLine string

		for i, s := range line {
			if i != 0 {
				sendLine += s
				sendLine += " "
			}
		}

		sendLine = strings.TrimSpace(sendLine)
		wsock(sendLine)

		// traps for regular shells and unwanted programs
	case "bash":
		fallthrough
	case "zsh":
		fallthrough
	case "sh":
		fallthrough
	case "fish":
		fmt.Println("External shell programs are disallowed")
		// --

	default:
		cmd := exec.Command(line[0], line[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		// Play ad:
		play_ad(PROGRAMMER)
		cmd.Run()

	}

	return 0
}

func splitLine(line string) []string {
	// TODO: argparse?
	return strings.Split(line, " ")
}
