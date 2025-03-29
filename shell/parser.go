package main

import (
	"github.com/inancgumus/screen"
	"os"
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
	case "zsh":
	case "sh":
	case "fish":
		// --
	}

	return 0
}

func splitLine(line string) []string {
	// TODO: argparse?
	return strings.Split(line, " ")
}
