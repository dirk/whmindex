package main

import (
	"fmt"
	"log"
	"os"
)

func printUsageAndExit() {
	fmt.Println("Usage: whmindex [compile|serve]")
	os.Exit(1)
}

func main() {
	if len(os.Args) != 2 {
		printUsageAndExit()
	}

	command := os.Args[1]
	var err error
	if command == "compile" {
		err = compile()
	} else if command == "serve" {
		err = serve()
	} else {
		printUsageAndExit()
	}
	if err != nil {
		log.Fatal(err)
	}
}
