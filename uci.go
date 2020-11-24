package main

import (
	"fmt"
)

func uci(frGUI chan string, tell func(text ...string)) {
	fmt.Println("Hello from uci!")
}

func mainTell(text ...string) {
	toGUI := ""
	for _, t := range text {
		toGUI += t
	}
	fmt.Println()
}
