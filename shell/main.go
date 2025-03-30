package main

import (
	"bufio"
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"log"
	"os"
	"os/signal"
	"os/user"
	"strings"
	"syscall"
)

var style = lipgloss.NewStyle().Foreground(lipgloss.Color("#D500D5"))

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
	}()

	inputCount := 0
	var coolDown = 5

	sockRead := make(chan []byte, 3)
	go func() {
		for {
			readData := rsock()
			if readData != nil {
				sockRead <- readData
				fmt.Println("\n" + string(<-sockRead))
				prompt(inputCount, coolDown)
			}
		}
	}()

	for {
		prompt(inputCount, coolDown)
		line := scanner()
		lines := splitLine(line)
		inputCount = parse(lines, inputCount, coolDown, sockRead)
		if inputCount > 5 {
			inputCount = 0
		}
	}
}

func prompt(inputCount int, coolDown int) {
	getwd, err := os.Getwd()
	if err != nil {
		return
	}
	homedir, err := os.UserHomeDir()
	if err != nil {
		return
	}

	current, err := user.Current()
	if err != nil {
		return
	}

	hostname, err := os.Hostname()
	if err != nil {
		return
	}

	remaningCommands := coolDown - inputCount
	if strings.Contains(getwd, homedir) {
		getwd = strings.Replace(getwd, homedir, "", 1)
		fmt.Print(style.Render(fmt.Sprintf(`╭─[%s@%s:~%s]
╰──────────── (%d Commands Remaning)$ `, current.Name, hostname, getwd, remaningCommands)))
	} else {
		fmt.Print(style.Render(fmt.Sprintf(`╭─[%s@%s:~%s]
╰─ (%d Commands Remaning)$ `, current.Name, hostname, getwd, remaningCommands)))
	}
}

func scanner() string {
	//line := ""
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	return line
}
