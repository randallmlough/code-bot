package main

import (
	"bufio"
	"os"
)

func waitForInput() string {
	var input string
	scanner := bufio.NewScanner(os.Stdin)
	for {
		// Scans a line from Stdin(Console)
		scanner.Scan()
		// Holds the string that scanned
		input = scanner.Text()
		if len(input) != 0 {
			return input
		} else {
			break
		}
	}
	return input
}
