package main

import (
	"fmt"
	"github.com/google/shlex"
	"github.com/inancgumus/screen"
	"os"
	"os/exec"
	"strings"
)

func parse(line []string, inputCount int, coolDown int, tagChan chan []byte) int {

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
		return inputCount
	case "pwd":
		getwd, err := os.Getwd()
		if err != nil {
			return 0
		}
		println(getwd)
		return inputCount
	case "exit":
		//err := os.Remove(socketPath)
		//if err != nil {
		//	os.Exit(55)
		//}
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
		return inputCount
	case "clear":
		screen.Clear()
		screen.MoveTopLeft()
		return inputCount
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
		return inputCount

		// traps for regular shells and unwanted programs
	case "bash":
		fallthrough
	case "zsh":
		fallthrough
	case "sh":
		fallthrough
	case "fish":
		fmt.Println("External shell programs are disallowed")
		return inputCount
		// --

	default:
		cmd := exec.Command(line[0], line[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		var fullLine string
		for _, str := range line {
			fullLine += str
		}
		//fmt.Println(getGroqCatagory(fullLine))
		// Play ad:
		if coolDown-inputCount == 1 {
			wsock(fullLine)
			tags := string(<-tagChan)
			play_ad(tags)
			fmt.Println(tags)

			//fmt.Println()
		}
		cmd.Run()

	}

	return inputCount + 1

}

func splitLine(line string) []string {
	split, err := shlex.Split(line)
	if err != nil {
		return nil
	}
	return split
}
