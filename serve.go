package main

import (
	"fmt"
	"net/http"
)

func handleSearch(res http.ResponseWriter, req *http.Request) {
	query := req.FormValue("query")
	if query == "" {
		res.Write([]byte(`{"message":"Missing query"}`))
	}
	res.Write([]byte(`{}`))
}

func serve() error {
	fmt.Printf("Building index...")
	err := buildIndex()
	if err != nil {
		fmt.Printf("\n")
		return err
	}
	fmt.Printf(" Done\n")

	fileServer := http.FileServer(http.Dir("_site"))
	http.Handle("/", fileServer)
	http.HandleFunc("/api/search.json", handleSearch)

	port := 3000
	fmt.Printf("Listening on %v...\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}
