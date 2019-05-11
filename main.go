package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: whmindex [compile|serve]")
		os.Exit(1)
	}

	command := os.Args[1]
	var err error
	if command == "compile" {
		err = compile()
	}
	if err != nil {
		log.Fatal(err)
	}
}
