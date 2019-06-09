package main

import (
	"encoding/json"
	"fmt"
	tmpl "html/template"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
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
var rootTemplate *tmpl.Template

func handleIndex(res http.ResponseWriter, _ *http.Request) error {
	return rootTemplate.Lookup("index.gohtml").Execute(res, index)
}

func handleEpisode(res http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	feed := vars["feed"]
	number, err := strconv.Atoi(vars["number"])
	if err != nil {
		return err
	}
	episode := index.FindEpisode(feed, number)
	if episode == nil {
		return fmt.Errorf("Episode not found (feed = %v, number = %v)", feed, number)
	}
	return rootTemplate.Lookup("episode.gohtml").Execute(res, episode)
}

func handleApiSearch(res http.ResponseWriter, req *http.Request) error {
	input := req.FormValue("query")
	if input == "" {
		return respondJson(res, http.StatusUnprocessableEntity, `{"message":"Missing query"}`)
	}
	query := parseQuery(input)
	result := executeSearch(index, query)
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

	fmt.Printf("Loading templates...")
	rootTemplate = tmpl.Must(tmpl.ParseGlob("templates/*.gohtml"))
	fmt.Printf(" Done\n")

	router := mux.NewRouter()
	router.HandleFunc("/api/search.json", handleError(handleApiSearch))
	router.HandleFunc("/", handleError(handleIndex))
	router.PathPrefix("/{feed:main}/{number:[0-9]+}").HandlerFunc(handleError(handleEpisode))
	fileServer := http.FileServer(http.Dir("static"))
	router.PathPrefix("/").Handler(fileServer)

	port := 3000
	fmt.Printf("Listening on %v...\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%v", port), router)
}
