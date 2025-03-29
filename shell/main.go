package main

import (
	"bufio"
	"log"
	"os"
)

func main() {
	for {
		print("# ")
		line := scanner()
		print(line)
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
