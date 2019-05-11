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
	files, err := filepath.Glob("data/transcripts/*.json")
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
	Results *TranscriptResult `json:"results"`
}

func (transcript *Transcript) CombinedTranscript() string {
	transcripts := make([]string, 0)
	for _, resultTranscript := range transcript.Results.Transcripts {
		transcripts = append(transcripts, resultTranscript.Transcript)
	}
	return strings.Join(transcripts, " ")
}

type TranscriptResult struct {
	Transcripts []*TranscriptResultTranscript `json:"transcripts"`
}

type TranscriptResultTranscript struct {
	Transcript string `json:"transcript"`
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

	outputFile := "data/episodes/" + strings.TrimSuffix(path.Base(file), ".json") + ".txt"
	return ioutil.WriteFile(outputFile, []byte(newlinedTranscript+"\n"), os.ModePerm)
}

func newlineSentences(transcript string) string {
	re := regexp.MustCompile(`[.?!] `)
	return re.ReplaceAllStringFunc(transcript, func(separator string) string {
		return strings.TrimSpace(separator) + "\n"
	})
}
