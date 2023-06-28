package main

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	green  = color.New(color.FgGreen).SprintFunc()
	blue   = color.New(color.FgBlue).SprintFunc()
	yellow = color.New(color.FgCyan).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
)

func printUserMessage(msg string) {
	fmt.Printf("%s %s\n", green("user:"), msg)
}
func printSystemMessage(msg string) {
	fmt.Printf("%s %s\n", blue("system:"), msg)
}
func printBotMessage(msg string) {
	fmt.Printf("%s %s\n", yellow("assistant:"), msg)
}
func printErrorMessage(msg string) {
	fmt.Printf("%s %s\n", red("ERROR:"), msg)
}
