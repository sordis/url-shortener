package main
//gocov:ignore

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if !shouldIgnore(line) {
			fmt.Println(line)
		}
	}
}

func shouldIgnore(line string) bool {
	// Игнорируем stdlib
	if strings.Contains(line, "/go/src/") {
		return true
	}


	ignorePhrases := []string{
		"log.Fatal",
		"os.Exit",
		"panic(",
	}

	for _, phrase := range ignorePhrases {
		if strings.Contains(line, phrase) {
			return true
		}
	}

	return false
}