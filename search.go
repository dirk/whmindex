package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"regexp"
	"strings"
)

// Structure of the `_data/index.yaml` file.
type DataIndex struct {
	Main []DataIndexEntry `yaml:"main"`
}

type DataIndexEntry struct {
	Number     int    `yaml:"number"`
	Title      string `yaml:"title"`
	Transcript string `yaml:"transcript"`
}

type Index struct {
	Episodes []IndexEpisode
}

type IndexEpisode struct {
	Feed   string
	Number int
	Lines  []IndexLine
}

type IndexLine struct {
	Contents string
	Words    map[string]int
}

func buildIndex() error {
	var dataIndex DataIndex
	contents, err := ioutil.ReadFile("_data/index.yaml")
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(contents, &dataIndex)
	if err != nil {
		return err
	}

	var index Index
	for _, dataIndexEntry := range dataIndex.Main {
		contents, err = ioutil.ReadFile("_data/transcripts/" + dataIndexEntry.Transcript)
		if err != nil {
			return err
		}

		episode := IndexEpisode{
			Feed:   "main",
			Number: dataIndexEntry.Number,
			Lines:  normalizeTranscriptForIndex(string(contents)),
		}
		index.Episodes = append(index.Episodes, episode)
	}

	// fmt.Printf("index: %#v\n", index)

	return nil
}

func normalizeTranscriptForIndex(transcript string) []IndexLine {
	lines := make([]IndexLine, 0)
	for _, line := range strings.Split(transcript, "\n") {
		normalizedLine := normalizeLineForIndex(line)
		if normalizedLine == "" {
			continue
		}

		words := make(map[string]int)
		for _, word := range strings.Split(normalizedLine, " ") {
			if count, exists := words[word]; exists {
				words[word] = count + 1
			} else {
				words[word] = 1
			}
		}

		lines = append(lines, IndexLine{
			Contents: normalizedLine,
			Words:    words,
		})
	}
	return lines
}

var punctuationRegexp = regexp.MustCompile(`[,.?!]( |$)`)

func normalizeLineForIndex(line string) string {
	downcased := strings.ToLower(line)
	withoutPunctuation := punctuationRegexp.ReplaceAllStringFunc(downcased, func(match string) string {
		if len(match) == 2 && match[1] == ' ' {
			return " "
		}
		return ""
	})
	return withoutPunctuation
}
