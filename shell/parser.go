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
		err := os.Chdir(line[1])
		if err != nil {
			return 0
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
	}
	return 0
}
