package main

import (
	"fmt"
	"net/http"
)

func serve() error {
	fileServer := http.FileServer(http.Dir("_site"))
	http.Handle("/", fileServer)

	port := 3000
	fmt.Printf("Listening on %v...\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}
