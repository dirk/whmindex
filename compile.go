package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

func compile() error {
	files, err := filepath.Glob("_data/sources/*.json")
	if err != nil {
		return err
	}
	for _, file := range files {
		err := compileFile(file)
		if err != nil {
			return err
		}
	}
	return nil
}

type Transcript struct {
	Results struct {
		Transcripts []struct {
			Transcript string `json:"transcript"`
		} `json:"transcripts"`
	} `json:"results"`
}

func (transcript *Transcript) CombinedTranscript() string {
	transcripts := make([]string, 0)
	for _, resultTranscript := range transcript.Results.Transcripts {
		transcripts = append(transcripts, resultTranscript.Transcript)
	}
	return strings.Join(transcripts, " ")
}

func compileFile(file string) error {
	var transcript Transcript
	contents, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(contents, &transcript)
	if err != nil {
		return err
	}

	combinedTranscript := transcript.CombinedTranscript()
	newlinedTranscript := newlineSentences(combinedTranscript)

	outputFile := "_data/transcripts/" + strings.TrimSuffix(path.Base(file), ".json") + ".txt"
	return ioutil.WriteFile(outputFile, []byte(newlinedTranscript+"\n"), os.ModePerm)
}

var newlineRegexp = regexp.MustCompile(`[.?!] `)

func newlineSentences(transcript string) string {
	return newlineRegexp.ReplaceAllStringFunc(transcript, func(separator string) string {
		return strings.TrimSpace(separator) + "\n"
	})
}
