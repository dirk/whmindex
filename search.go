package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"regexp"
	"sort"
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
	Episodes []*IndexEpisode
}

func (index *Index) FindEpisode(feed string, number int) *IndexEpisode {
	for _, episode := range index.Episodes {
		if episode.Feed == feed && episode.Number == number {
			return episode
		}
	}
	return nil
}

type IndexEpisode struct {
	Feed   string
	Number int
	Title  string
	Lines  []IndexLine
}

type IndexLine struct {
	Contents string
	Words    map[string]int
}

func buildIndex() (*Index, error) {
	var dataIndex DataIndex
	contents, err := ioutil.ReadFile("_data/index.yaml")
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(contents, &dataIndex)
	if err != nil {
		return nil, err
	}

	var index Index
	for _, dataIndexEntry := range dataIndex.Main {
		contents, err = ioutil.ReadFile("_data/transcripts/" + dataIndexEntry.Transcript)
		if err != nil {
			return nil, err
		}

		episode := &IndexEpisode{
			Feed:   "main",
			Number: dataIndexEntry.Number,
			Title:  dataIndexEntry.Title,
			Lines:  normalizeTranscriptForIndex(string(contents)),
		}
		index.Episodes = append(index.Episodes, episode)
	}

	// fmt.Printf("index: %#v\n", index)

	return &index, nil
}

var spacesRegexp = regexp.MustCompile(`\s+`)

func normalizeTranscriptForIndex(transcript string) []IndexLine {
	lines := make([]IndexLine, 0)
	for _, line := range strings.Split(transcript, "\n") {
		normalizedLine := normalizeLineForIndex(line)
		if normalizedLine == "" {
			continue
		}

		words := make(map[string]int)
		for _, word := range spacesRegexp.Split(normalizedLine, -1) {
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

type Query struct {
	// Single word terms (can use the `Words` frequency map in `Line`).
	Words []string
	// Phrases that must be substring searched in each line.
	Phrases []string
}

var singleQuotedPhraseRegexp = regexp.MustCompile(`'[^']*'`)
var doubleQuotedPhraseRegexp = regexp.MustCompile(`"[^"]*"`)
var allNonWordOrSpaceRegexp = regexp.MustCompile(`[^A-Za-z0-9 ]`)

func extractQuotedPhrases(input string, quotedPhraseRegexp *regexp.Regexp) (string, []string) {
	phrases := make([]string, 0)
	replaced := quotedPhraseRegexp.ReplaceAllStringFunc(input, func(match string) string {
		phrases = append(phrases, match[1:len(match)-1])
		return ""
	})
	return replaced, phrases
}

func parseQuery(input string) *Query {
	// First normalize everything to lowercase.
	input = strings.ToLower(input)

	query := &Query{
		Words:   make([]string, 0),
		Phrases: make([]string, 0),
	}

	// Extract single- and double-quoted sequences into phrases.
	var phrases []string
	input, phrases = extractQuotedPhrases(input, singleQuotedPhraseRegexp)
	query.Phrases = append(query.Phrases, phrases...)
	input, phrases = extractQuotedPhrases(input, doubleQuotedPhraseRegexp)
	query.Phrases = append(query.Phrases, phrases...)

	// Remove all characters besides letters, numbers, and spaces.
	input = allNonWordOrSpaceRegexp.ReplaceAllString(input, "")
	input = strings.TrimSpace(input)
	// Split on 1+ spaces.
	query.Words = spacesRegexp.Split(input, -1)

	return query
}

type Result struct {
	Matches []*Match `json:"matches"`
}

type Match struct {
	Episode *MatchEpisode `json:"episode"`
	Score   int           `json:"score"`
}

type MatchEpisode struct {
	Feed   string `json:"feed"`
	Number int    `json:"number"`
	Title  string `json:"title"`
}

func NewMatchEpisode(episode *IndexEpisode) *MatchEpisode {
	return &MatchEpisode{
		Feed:   episode.Feed,
		Number: episode.Number,
		Title:  episode.Title,
	}
}

func executeSearch(index *Index, query *Query) Result {
	matches := make([]*Match, 0)
	for _, episode := range index.Episodes {
		score := 0
		for _, line := range episode.Lines {
			score += scoreLine(&line, query)
		}
		if score > 0 {
			matches = append(matches, &Match{
				Episode: NewMatchEpisode(episode),
				Score:   score,
			})
		}
	}
	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Score > matches[j].Score
	})
	return Result{Matches: matches}
}

func scoreLine(line *IndexLine, query *Query) int {
	score := 0
	for _, word := range query.Words {
		if count, exists := line.Words[word]; exists {
			score += count
		}
	}
	// TODO: Check phrases if any.
	return score
}
