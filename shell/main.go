package main

import (
	"bufio"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
	}()

	for {
		print("# ")
		line := scanner()
		lines := strings.Split(line, " ")
		parse(lines)
		//print()
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
