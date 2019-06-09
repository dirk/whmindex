package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Returns type of `error` but will actually always return `nil`. That way this
// can be used like `return respondJson(...)`.
func respondJson(res http.ResponseWriter, code int, body string) error {
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(code)
	_, err := res.Write([]byte(body))
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error writing response: %v\n", err)
	}
	return nil
}

var index *Index

func handleSearch(res http.ResponseWriter, req *http.Request) error {
	input := req.FormValue("query")
	if input == "" {
		return respondJson(res, http.StatusUnprocessableEntity, `{"message":"Missing query"}`)
	}
	query := parseQuery(input)
	// fmt.Printf("query: %#v\n", query)

	result := executeSearch(index, query)
	fmt.Printf("result: %#v\n", result)

	body, err := json.Marshal(result)
	if err != nil {
		return err
	}
	return respondJson(res, http.StatusOK, string(body))
}

func handleError(handler func(http.ResponseWriter, *http.Request) error) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		err := handler(res, req)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error handling request: %#v\n", err)
			_ = respondJson(res, http.StatusInternalServerError, `{"message":"Internal server error"}`)
		}
	}
}

func serve() error {
	fmt.Printf("Building index...")
	var err error
	index, err = buildIndex()
	if err != nil {
		fmt.Printf("\n")
		return err
	}
	fmt.Printf(" Done\n")

	fileServer := http.FileServer(http.Dir("_site"))
	http.Handle("/", fileServer)
	http.HandleFunc("/api/search.json", handleError(handleSearch))

	port := 3000
	fmt.Printf("Listening on %v...\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}
