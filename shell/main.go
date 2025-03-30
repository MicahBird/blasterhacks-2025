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

const socketPath = "/tmp/unix.sock"

var style = lipgloss.NewStyle().Foreground(lipgloss.Color("#D500D5"))

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
	}()

	sockRead := make(chan []byte, 3)
	go func() {
		for {
			readData := rsock()
			if readData != nil {
				sockRead <- readData
				fmt.Println("\n" + string(<-sockRead))
				prompt()
			}
		}
	}()

	inputCount := 0

	for {
		prompt()
		line := scanner()
		lines := splitLine(line)
		inputCount += 1
		parse(lines, inputCount)
		//print()
	}
}

func prompt() {
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

	if strings.Contains(getwd, homedir) {
		getwd = strings.Replace(getwd, homedir, "", 1)
		fmt.Print(style.Render(fmt.Sprintf("[%s@%s:~%s]$ ", current.Name, hostname, getwd)))
	} else {
		fmt.Print(style.Render(fmt.Sprintf("[%s@%s:%s]$ ", current.Name, hostname, getwd)))
	}
}

func scanner() string {
	//line := ""
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
		os.Remove(socketPath)
	}

	return line
}
